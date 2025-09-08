package application

import (
	"context"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
)

type TeacherQueryService interface {
	GetTeacherList(ctx context.Context, query course.TeacherListQuery) ([]vo.TeacherListItemVO, error)
	GetTeacherFilter(ctx context.Context) (*course.TeacherFilter, error)
}
