package query

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

// 设计说明（Query 层为何只用 GORM 而非 gorm/gen）
// 本文件中的查询服务遵循 CQRS 的“读侧只为展示塑形”的目标：
// - 直接用 GORM Table/Joins/Where/Select/Scan 构造包含 MAX、EXISTS 等子查询与聚合的 SQL，
//   可读性更强、拼装更灵活，便于实现“最新学期”“按存在学期过滤”等复杂筛选。
// - 只选择必要列并扫描到轻量 VO/临时结构，避免实体预加载带来的冗余 IO 与字段。
// - 降低对代码生成器的耦合，读侧需求频繁变动时无需反复生成代码。
// - 写侧/仓储仍可使用 gen；读侧为适配跨表聚合与视图化塑形采用更自由的 GORM。
// 详细说明见 README.md 中《为什么查询层只使用 GORM，而不直接使用 gorm/gen？》一节。

type CourseQueryService interface {
	GetCourseList(ctx context.Context, query course.CourseListQuery) ([]vo.CourseListItemVO, error)
	GetCourseDetail(ctx context.Context, courseID shared.IDType) (*vo.CourseDetailVO, error)
	GetCourseFilter(ctx context.Context) (*course.CourseFilter, error)
}

type courseQueryService struct {
	db *gorm.DB
}

func (s *courseQueryService) GetCourseList(ctx context.Context, q course.CourseListQuery) ([]vo.CourseListItemVO, error) {
	if s.db == nil {
		return nil, errors.New("db not initialized")
	}

	db := s.db.WithContext(ctx)

	// Defaults for pagination
	limit := int(q.PageSize)
	if limit <= 0 {
		limit = 20
	}
	page := int(q.Page)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// Base query joining latest offering and main teacher
	base := db.Model(&entity.Course{}).
		Joins("LastOffering").
		Joins("MainTeacher")

	// Filters
	if q.Code != "" {
		base = base.Where("code = ?", q.Code)
	}
	if q.MainTeacherID > 0 {
		base = base.Where("main_teacher_id = ?", q.MainTeacherID)
	}
	if len(q.Credits) > 0 {
		base = base.Where("credit IN ?", q.Credits)
	}
	// Filter courses that have offerings in any of the specified semesters
	if len(q.Semesters) > 0 {
		existsSub := db.Table("course_offering co").Select("1").
			Where("co.course_id = c.id AND co.semester IN ?", q.Semesters)
		base = base.Where("EXISTS (?)", existsSub)
	}
	// Departments/Categories filter on latest offering only
	if len(q.Departments) > 0 {
		base = base.Where("last_offering.department IN ?", q.Departments)
	}
	if len(q.Categories) > 0 {
		base = base.Joins("JOIN course_offering_category coc ON coc.course_offering_id = last_offering_id").
			Where("coc.category IN ?", q.Categories)
	}

	// Select distinct to avoid duplication due to category joins
	rows := make([]entity.Course, 0)
	err := base.
		Limit(limit).Offset(offset).
		Order("code ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []vo.CourseListItemVO{}, nil
	}

	// Build result
	result := make([]vo.CourseListItemVO, 0, len(rows))
	for _, r := range rows {

		catagories := make([]string, 0)
		for _, ct := range r.LastOffering.Categories {
			catagories = append(catagories, ct.Category)
		}

		item := vo.CourseListItemVO{
			ID:     r.ID,
			Code:   r.Code,
			Name:   r.Name,
			Credit: r.Credit,
			MainTeacher: vo.TeacherInCourseVO{
				ID:         r.MainTeacherID,
				Name:       r.MainTeacher.Name,
				Department: r.MainTeacher.Department,
			},
			LatestOffering: vo.OfferingInfoVO{
				Categories: catagories,
				Department: r.LastOffering.Department,
				Language:   r.LastOffering.Language,
				Semester:   r.LastOffering.Semester,
			},
			RatingInfo: vo.RatingVO{
				Average: r.ReviewAvg,
				Count:   r.ReviewCount,
			},
		}
		result = append(result, item)
	}

	return result, nil
}

func (s *courseQueryService) GetCourseDetail(ctx context.Context, courseID shared.IDType) (*vo.CourseDetailVO, error) {
	c := entity.Course{}
	if err := s.db.WithContext(ctx).
		Model(&entity.Course{}).
		Joins("MainTeacher").
		Joins("LastOffering").
		Preload("Offerings.TeacherGroup.Teacher").
		Preload("Offerings.Categories").
		Where("course.id = ?", courseID).Take(&c).Error; err != nil {
		return nil, err
	}

	offeringVOs := make([]vo.CourseOfferingVO, 0)
	for _, o := range c.Offerings {
		categories := make([]string, 0)
		for _, ct := range o.Categories {
			categories = append(categories, ct.Category)
		}
		teacherGroup := make([]vo.TeacherInCourseVO, 0)
		for _, t := range o.TeacherGroup {
			teacherGroup = append(teacherGroup, vo.TeacherInCourseVO{
				ID:         t.Teacher.ID,
				Name:       t.Teacher.Name,
				Department: t.Teacher.Department,
			})
		}
		offeringVOs = append(offeringVOs, vo.CourseOfferingVO{
			Semester:     o.Semester,
			TeacherGroup: teacherGroup,
			Categories:   categories,
			Department:   o.Department,
			Language:     o.Language,
		})
	}

	return &vo.CourseDetailVO{
		ID:     c.ID,
		Code:   c.Code,
		Name:   c.Name,
		Credit: c.Credit,
		MainTeacher: vo.TeacherInCourseVO{
			ID:         c.MainTeacher.ID,
			Name:       c.MainTeacher.Name,
			Department: c.MainTeacher.Department,
		},
		Offering: offeringVOs,
	}, nil
}

func (s *courseQueryService) GetCourseFilter(ctx context.Context) (*course.CourseFilter, error) {
	if s.db == nil {
		return nil, errors.New("db not initialized")
	}

	// latest offering per course (string compare for semester)
	latestSub := s.db.WithContext(ctx).Model(&entity.CourseOffering{}).
		Select("course_id, MAX(semester) as latest_semester").
		Group("course_id")

	getBase := func() *gorm.DB {
		return s.db.WithContext(ctx).Table("course as c").
			Joins("JOIN (?) latest ON latest.course_id = c.id", latestSub).
			Joins("JOIN course_offering lo ON lo.course_id = c.id AND lo.semester = latest.latest_semester")
	}

	filter := &course.CourseFilter{
		Departments: make([]course.FilterAggregateItem, 0),
		Credits:     make([]course.FilterAggregateItem, 0),
		Semesters:   make([]course.FilterAggregateItem, 0),
		Categories:  make([]course.FilterAggregateItem, 0),
	}

	if err := getBase().
		Select("c.credit as value, COUNT(DISTINCT c.id) as count").
		Group("c.credit").
		Order("c.credit ASC").
		Scan(&filter.Credits).Error; err != nil {
		return nil, err
	}

	// Departments from latest offering only

	if err := getBase().
		Select("lo.department as value, COUNT(DISTINCT c.id) as count").
		Group("lo.department").
		Order("lo.department ASC").
		Scan(&filter.Departments).Error; err != nil {
		return nil, err
	}

	// Semesters: latest semester value per course
	if err := getBase().
		Select("latest.latest_semester as value, COUNT(DISTINCT c.id) as count").
		Group("latest.latest_semester").
		Order("latest.latest_semester DESC").
		Scan(&filter.Semesters).Error; err != nil {
		return nil, err
	}

	// Categories from latest offering only
	if err := getBase().
		Joins("JOIN course_offering_category coc ON coc.course_offering_id = lo.id").
		Select("coc.category as value, COUNT(DISTINCT c.id) as count").
		Group("coc.category").
		Order("coc.category ASC").
		Scan(&filter.Categories).Error; err != nil {
		return nil, err
	}
	return filter, nil
}

func NewCourseQueryService(db *gorm.DB) CourseQueryService {
	return &courseQueryService{db: db}
}
