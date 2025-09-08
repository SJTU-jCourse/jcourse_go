package application

import (
	"context"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
)

type CourseQueryService interface {
	GetCourseList(ctx context.Context, query course.CourseListQuery) ([]vo.CourseListItemVO, error)
	GetCourseDetail(ctx context.Context, courseID shared.IDType) (*vo.CourseDetailVO, error)
	GetCourseFilter(ctx context.Context) (*course.CourseFilter, error)
}
