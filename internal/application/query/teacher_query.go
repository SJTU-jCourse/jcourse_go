package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
)

type TeacherQueryService interface {
	GetTeacherList(ctx context.Context, query course.TeacherListQuery) ([]vo.TeacherListItemVO, error)
	GetTeacherFilter(ctx context.Context) (*course.TeacherFilter, error)
}

type teacherQueryService struct {
	db *gorm.DB
}

func (t *teacherQueryService) GetTeacherList(ctx context.Context, query course.TeacherListQuery) ([]vo.TeacherListItemVO, error) {
	// TODO implement me
	panic("implement me")
}

func (t *teacherQueryService) GetTeacherFilter(ctx context.Context) (*course.TeacherFilter, error) {
	// TODO implement me
	panic("implement me")
}

func NewTeacherQueryService(db *gorm.DB) TeacherQueryService {
	return &teacherQueryService{db: db}
}
