package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type CourseQueryService interface {
	GetCourseList(ctx context.Context, query course.CourseListQuery) ([]vo.CourseListItemVO, error)
	GetCourseDetail(ctx context.Context, courseID shared.IDType) (*vo.CourseDetailVO, error)
	GetCourseFilter(ctx context.Context) (*course.CourseFilter, error)
}

type courseQueryService struct {
	db *gorm.DB
}

func (s *courseQueryService) GetCourseList(ctx context.Context, query course.CourseListQuery) ([]vo.CourseListItemVO, error) {

	latestOfferingSubQuery := s.db.Model(&entity.CourseOffering{}).Select("course_id, max(semester) as latest_semester").Group("course_id")

	q := s.db.Model("courses as c").
		Joins("LEFT JOIN course_offering as co ON c.id = co.course_id").
		Joins("LEFT JOIN teacher as t ON c.main_teacher_id = t.id").
		Joins("JOIN (?) as latest_co ON c.id = latest_co.course_id", latestOfferingSubQuery)

	if query.Code != "" {
		q = q.Where("c.code = ?", query.Code)
	}
	if len(query.Categories) > 0 {
		q = q.Joins("LEFT JOIN course_offering_categories as coc ON latest_co.id = coc.course_offering_id").
			Where("coc.category_id IN (?)", query.Categories)
	}
	if len(query.Credits) > 0 {
		q = q.Where("c.credits in (?)", query.Credits)
	}
	if len(query.Departments) > 0 {
		q = q.Where("latest_co.department_id IN (?)", query.Departments)
	}
	if query.MainTeacherID != 0 {
		q = q.Where("c.main_teacher_id = ?", query.MainTeacherID)
	}

	return nil, nil
}

func (s *courseQueryService) GetCourseDetail(ctx context.Context, courseID shared.IDType) (*vo.CourseDetailVO, error) {
	return nil, nil
}

func (s *courseQueryService) GetCourseFilter(ctx context.Context) (*course.CourseFilter, error) {
	// TODO implement me
	panic("implement me")
}

func NewCourseQueryService(db *gorm.DB) CourseQueryService {
	return &courseQueryService{db: db}
}
