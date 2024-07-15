package repository

import (
	"context"
	"errors"

	"jcourse_go/dal"
	"jcourse_go/model/po"

	"gorm.io/gorm"
)

type ITeacherQuery interface {
	GetTeacher(ctx context.Context, opts ...DBOption) (*po.TeacherPO, error)
	GetTeacherList(ctx context.Context, opts ...DBOption) ([]po.TeacherPO, error)
	WithID(id int64) DBOption
	WithCode(code string) DBOption
	WithName(name string) DBOption
	WithPinyin(pinyin string) DBOption
	WithPinyinAbbr(pinyin string) DBOption
	WithDepartment(department string) DBOption
	WithTitle(title string) DBOption
	WithPicture(picture string) DBOption
	WithProfileURL(profileURL string) DBOption
}

type TeacherQuery struct {
	db *gorm.DB
}

func NewTeacherQuery() ITeacherQuery {
	return &TeacherQuery{
		db: dal.GetDBClient(),
	}
}

func (q *TeacherQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := q.db.WithContext(ctx).Model(po.TeacherPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (q *TeacherQuery) GetTeacher(ctx context.Context, opts ...DBOption) (*po.TeacherPO, error) {
	db := q.optionDB(ctx, opts...)
	teacher := po.TeacherPO{}
	result := db.First(&teacher)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &teacher, nil
}

func (q *TeacherQuery) GetTeacherList(ctx context.Context, opts ...DBOption) ([]po.TeacherPO, error) {
	db := q.optionDB(ctx, opts...)
	teacherPOs := make([]po.TeacherPO, 0)
	result := db.Find(&teacherPOs)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return teacherPOs, nil
}

func (q *TeacherQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (q *TeacherQuery) WithCode(code string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("code = ?", code)
	}
}

func (q *TeacherQuery) WithName(name string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func (q *TeacherQuery) WithPinyin(pinyin string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("pinyin = ?", pinyin)
	}
}

func (q *TeacherQuery) WithPinyinAbbr(pinyin string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("pinyin_abbr = ?", pinyin)
	}
}

func (q *TeacherQuery) WithDepartment(department string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("department = ?", department)
	}
}

func (q *TeacherQuery) WithTitle(title string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("title = ?", title)
	}
}

func (q *TeacherQuery) WithPicture(picture string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("picture = ?", picture)
	}
}

func (q *TeacherQuery) WithProfileURL(profileURL string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("profile_url = ?", profileURL)
	}
}

type TeacherCourseQuery struct {
	db *gorm.DB
}

type ITeacherCourseQuery interface {
	GetTeacherCourseList(ctx context.Context, opts ...DBOption) ([]po.TeacherCoursePO, error)
	GetTeacherBaseCourseList(ctx context.Context, teacherID int64) ([]po.BaseCoursePO, error)
	GetTeacherCourseIDs(ctx context.Context, opts ...DBOption) ([]int64, error)
	optionDB(ctx context.Context, opts ...DBOption) *gorm.DB
	WithTeacherID(teacherID int64) DBOption
	WithCourseID(courseID int64) DBOption
	WithCourseIDs(courseIDs []int64) DBOption
}

func NewTeacherCourseQuery() ITeacherCourseQuery {
	return &TeacherCourseQuery{
		db: dal.GetDBClient(),
	}
}

func (t *TeacherCourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := t.db.Model(&po.TeacherCoursePO{}).WithContext(ctx)
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (t *TeacherCourseQuery) GetTeacherCourseList(ctx context.Context, opts ...DBOption) ([]po.TeacherCoursePO, error) {
	db := t.optionDB(ctx, opts...)
	teacherCourses := make([]po.TeacherCoursePO, 0)
	result := db.Debug().Find(&teacherCourses)
	if result.Error != nil {
		return nil, result.Error
	}
	return teacherCourses, nil
}

func (t *TeacherCourseQuery) GetTeacherBaseCourseList(ctx context.Context, teacherID int64) ([]po.BaseCoursePO, error) {
	c := NewBaseCourseQuery()
	opt := t.WithTeacherID(teacherID)
	courseIDs := make([]int64, 0)
	t.optionDB(ctx, opt).Find(&courseIDs)
	copt := c.WithIDs(courseIDs)
	courses, err := c.GetBaseCourseList(ctx, copt)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func (t *TeacherCourseQuery) GetTeacherCourseIDs(ctx context.Context, opts ...DBOption) ([]int64, error) {
	TeacherCourses := make([]int64, 0)
	db := t.optionDB(ctx, opts...)
	result := db.Pluck("teacher_id", &TeacherCourses)
	if result.Error != nil {
		return nil, result.Error
	}
	return TeacherCourses, nil
}

func (t *TeacherCourseQuery) WithTeacherID(teacherID int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("teacher_id = ?", teacherID)
	}
}

func (t *TeacherCourseQuery) WithCourseID(courseID int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id = ?", courseID)
	}
}

func (t *TeacherCourseQuery) WithCourseIDs(courseIDs []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id IN ?", courseIDs).
			Group("teacher_id").
			Having("count(DISTINCT course_id) = ?", len(courseIDs))
	}
}
