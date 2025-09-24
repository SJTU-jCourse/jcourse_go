package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type CourseRepository struct {
	db *gorm.DB
}

func (r *CourseRepository) Get(ctx context.Context, id shared.IDType) (*course.Course, error) {
	e := &entity.Course{}
	if err := r.db.Model(&entity.Course{}).Where("id = ?", id).First(&e).Error; err != nil {
		return nil, err
	}
	d := newCourseDomainFromEntity(e)
	return &d, nil
}

func NewCourseRepository(db *gorm.DB) course.CourseRepository {
	return &CourseRepository{db: db}
}

func newCourseDomainFromEntity(c *entity.Course) course.Course {
	return course.Course{
		ID:            shared.IDType(c.ID),
		Code:          c.Code,
		Name:          c.Name,
		Credit:        c.Credit,
		MainTeacherID: shared.IDType(c.MainTeacherID),
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}
