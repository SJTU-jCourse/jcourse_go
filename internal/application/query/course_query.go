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

	// Subquery: latest semester per course
	latestSub := db.Table("course_offering as co").
		Select("co.course_id, MAX(co.semester) as latest_semester").
		Group("co.course_id")

	// Base query joining latest offering and main teacher
	base := db.Table("course as c").
		Joins("JOIN (?) latest ON latest.course_id = c.id", latestSub).
		Joins("JOIN course_offering lo ON lo.course_id = c.id AND lo.semester = latest.latest_semester").
		Joins("JOIN teacher t ON t.id = c.main_teacher_id")

	// Filters
	if q.Code != "" {
		base = base.Where("c.code = ?", q.Code)
	}
	if q.MainTeacherID > 0 {
		base = base.Where("c.main_teacher_id = ?", q.MainTeacherID)
	}
	if len(q.Credits) > 0 {
		base = base.Where("c.credit IN ?", q.Credits)
	}
	// Filter courses that have offerings in any of the specified semesters
	if len(q.Semesters) > 0 {
		existsSub := db.Table("course_offering co2").Select("1").
			Where("co2.course_id = c.id AND co2.semester IN ?", q.Semesters)
		base = base.Where("EXISTS (?)", existsSub)
	}
	// Departments/Categories filter on latest offering only
	if len(q.Departments) > 0 {
		base = base.Where("lo.department IN ?", q.Departments)
	}
	if len(q.Categories) > 0 {
		base = base.Joins("JOIN course_offering_category coc ON coc.course_offering_id = lo.id").
			Where("coc.category IN ?", q.Categories)
	}

	// Select distinct to avoid duplication due to category joins
	rows := make([]struct {
		ID                 int64   `json:"id"`
		Code               string  `json:"code"`
		Name               string  `json:"name"`
		Credit             float64 `json:"credit"`
		MainTeacherID      int64   `json:"main_teacher_id"`
		TeacherName        string  `json:"teacher_name"`
		TeacherDepartment  string  `json:"teacher_department"`
		OfferingID         int64   `json:"offering_id"`
		OfferingDepartment string  `json:"offering_department"`
		OfferingLanguage   string  `json:"offering_language"`
		LatestSemester     string  `json:"latest_semester"`
	}, 0)

	err := base.Select("DISTINCT " +
		"c.id as id, " +
		"c.code, " +
		"c.name, " +
		"c.credit, " +
		"c.main_teacher_id, " +
		"t.name as teacher_name, " +
		"t.department as teacher_department, " +
		"lo.id as offering_id, " +
		"lo.department as offering_department, " +
		"lo.language as offering_language, " +
		"latest.latest_semester").
		Limit(limit).Offset(offset).
		Order("c.code ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return []vo.CourseListItemVO{}, nil
	}

	// Collect offering IDs and course IDs
	offeringIDs := make([]int64, 0, len(rows))
	courseIDs := make([]int64, 0, len(rows))
	for _, r := range rows {
		offeringIDs = append(offeringIDs, r.OfferingID)
		courseIDs = append(courseIDs, r.ID)
	}

	// Fetch categories for the latest offerings
	catRows := make([]struct {
		CourseOfferingID int64  `json:"course_offering_id"`
		Category         string `json:"category"`
	}, 0)
	catMap := make(map[int64][]string)
	if len(offeringIDs) > 0 {
		err = db.Table("course_offering_category").
			Select("course_offering_id, category").
			Where("course_offering_id IN ?", offeringIDs).
			Scan(&catRows).Error
		if err != nil {
			return nil, err
		}
		for _, cr := range catRows {
			catMap[cr.CourseOfferingID] = append(catMap[cr.CourseOfferingID], cr.Category)
		}
	}

	// Compute rating from reviews
	ratingRows := make([]struct {
		CourseID int64   `json:"course_id"`
		Count    int64   `json:"count"`
		Average  float64 `json:"average"`
	}, 0)
	ratingMap := make(map[int64]vo.RatingVO)
	err = db.Table("review").
		Select("course_id, COUNT(*) as count, AVG(rating) as average").
		Where("course_id IN ?", courseIDs).
		Group("course_id").
		Scan(&ratingRows).Error
	if err != nil {
		return nil, err
	}
	for _, rr := range ratingRows {
		ratingMap[rr.CourseID] = vo.RatingVO{Average: rr.Average, Count: rr.Count}
	}

	// Build result
	result := make([]vo.CourseListItemVO, 0, len(rows))
	for _, r := range rows {
		item := vo.CourseListItemVO{
			ID:     r.ID,
			Code:   r.Code,
			Name:   r.Name,
			Credit: r.Credit,
			MainTeacher: vo.TeacherInCourseVO{
				ID:         r.MainTeacherID,
				Name:       r.TeacherName,
				Department: r.TeacherDepartment,
			},
			LatestOffering: vo.OfferingInfoVO{
				Categories: catMap[r.OfferingID],
				Department: r.OfferingDepartment,
				Language:   r.OfferingLanguage,
			},
			RatingInfo: ratingMap[r.ID],
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
