package repository

import (
	"context"

	"jcourse_go/model/domain"
)

type IBaseCourseQuery interface {
	GetBaseCourse(ctx context.Context, opts ...DBOption) (*domain.BaseCourse, error)
	GetBaseCourseList(ctx context.Context, opts ...DBOption) ([]domain.BaseCourse, error)
	WithCode(code string) DBOption
	WithName(name string) DBOption
	WithCredit(credit float64) DBOption
}

type ICourseQuery interface {
	GetCourse(ctx context.Context, opts ...DBOption) (*domain.Course, error)
	GetCourseList(ctx context.Context, opts ...DBOption) ([]domain.Course, error)
	WithID(id int64) DBOption
	WithBaseCourseID(id int64) DBOption
	WithCode(code string) DBOption
	WithName(name string) DBOption
	WithCredit(credit float64) DBOption
	WithCategory(category string) DBOption
	WithMainTeacherName(name string) DBOption
	WithMainTeacherID(id int64) DBOption
}

type IOfferedCourseQuery interface {
	GetOfferedCourse(ctx context.Context, opts ...DBOption) (*domain.OfferedCourse, error)
	GetOfferedCourseList(ctx context.Context, opts ...DBOption) ([]domain.OfferedCourse, error)
	WithID(id int64) DBOption
	WithCourseID(id int64) DBOption
	WithMainTeacherName(name string) DBOption
	WithMainTeacherID(id int64) DBOption
	WithSemester(semester string) DBOption
	WithDepartment(department string) DBOption
	WithLanguage(language string) DBOption
}
