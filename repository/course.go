package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"jcourse_go/dal"
	"jcourse_go/model/po"
)

type IBaseCourseQuery interface {
	GetBaseCourse(ctx context.Context, opts ...DBOption) (*po.BaseCoursePO, error)
	GetBaseCourseList(ctx context.Context, opts ...DBOption) ([]po.BaseCoursePO, error)
	WithCode(code string) DBOption
	WithName(name string) DBOption
	WithCredit(credit float64) DBOption
}

type BaseCourseQuery struct {
	db *gorm.DB
}

func (b *BaseCourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := b.db.WithContext(ctx).Model(po.BaseCoursePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (b *BaseCourseQuery) GetBaseCourse(ctx context.Context, opts ...DBOption) (*po.BaseCoursePO, error) {
	db := b.optionDB(ctx, opts...)
	course := po.BaseCoursePO{}
	result := db.Debug().First(&course)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (b *BaseCourseQuery) GetBaseCourseList(ctx context.Context, opts ...DBOption) ([]po.BaseCoursePO, error) {
	db := b.optionDB(ctx, opts...)
	coursePOs := make([]po.BaseCoursePO, 0)
	result := db.Debug().Find(&coursePOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return coursePOs, nil
}

func (b *BaseCourseQuery) WithCode(code string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func (b *BaseCourseQuery) WithName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (b *BaseCourseQuery) WithCredit(credit float64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("credit = ?", credit)
	}
}

func NewBaseCourseQuery() IBaseCourseQuery {
	return &BaseCourseQuery{db: dal.GetDBClient()}
}

type ICourseQuery interface {
	GetCourse(ctx context.Context, opts ...DBOption) (*po.CoursePO, error)
	GetCourseList(ctx context.Context, opts ...DBOption) ([]po.CoursePO, error)
	WithID(id int64) DBOption
	WithCode(code string) DBOption
	WithName(name string) DBOption
	WithCredit(credit float64) DBOption
	WithCategory(category string) DBOption
	WithMainTeacherName(name string) DBOption
	WithMainTeacherID(id int64) DBOption
}

type CourseQuery struct {
	db *gorm.DB
}

func (c *CourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := c.db.WithContext(ctx).Model(po.CoursePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (c *CourseQuery) GetCourse(ctx context.Context, opts ...DBOption) (*po.CoursePO, error) {
	db := c.optionDB(ctx, opts...)
	course := po.CoursePO{}
	result := db.Debug().First(&course)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (c *CourseQuery) GetCourseList(ctx context.Context, opts ...DBOption) ([]po.CoursePO, error) {
	db := c.optionDB(ctx, opts...)
	coursePOs := make([]po.CoursePO, 0)
	result := db.Debug().Find(&coursePOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return coursePOs, nil
}

func (c *CourseQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (c *CourseQuery) WithCode(code string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func (c *CourseQuery) WithName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (c *CourseQuery) WithCredit(credit float64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("credit = ?", credit)
	}
}

func (c *CourseQuery) WithCategory(category string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("inner join course_categories on course_categories.course_id = courses.id").Where("category = ?", category)
	}
}

func (c *CourseQuery) WithMainTeacherName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("main_teacher_name = ?", name)
	}
}

func (c *CourseQuery) WithMainTeacherID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("main_teacher_id = ?", id)
	}
}

func NewCourseQuery() ICourseQuery {
	return &CourseQuery{db: dal.GetDBClient()}
}

type IOfferedCourseQuery interface {
	GetOfferedCourse(ctx context.Context, opts ...DBOption) (*po.OfferedCoursePO, error)
	GetOfferedCourseList(ctx context.Context, opts ...DBOption) ([]po.OfferedCoursePO, error)
	WithID(id int64) DBOption
	WithCourseID(id int64) DBOption
	WithMainTeacherID(id int64) DBOption
	WithSemester(semester string) DBOption
}

type OfferedCourseQuery struct {
	db *gorm.DB
}

func (o *OfferedCourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := o.db.WithContext(ctx).Model(po.OfferedCoursePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (o *OfferedCourseQuery) GetOfferedCourse(ctx context.Context, opts ...DBOption) (*po.OfferedCoursePO, error) {
	db := o.optionDB(ctx, opts...)
	course := po.OfferedCoursePO{}
	result := db.Debug().First(&course)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (o *OfferedCourseQuery) GetOfferedCourseList(ctx context.Context, opts ...DBOption) ([]po.OfferedCoursePO, error) {
	db := o.optionDB(ctx, opts...)
	coursePOs := make([]po.OfferedCoursePO, 0)
	result := db.Debug().Find(&coursePOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return coursePOs, nil
}

func (o *OfferedCourseQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (o *OfferedCourseQuery) WithCourseID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id = ?", id)
	}
}

func (o *OfferedCourseQuery) WithMainTeacherID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("main_teacher_id = ?", id)
	}
}

func (o *OfferedCourseQuery) WithSemester(semester string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("semester = ?", semester)
	}
}

func NewOfferedCourseQuery() IOfferedCourseQuery {
	return &OfferedCourseQuery{db: dal.GetDBClient()}
}
