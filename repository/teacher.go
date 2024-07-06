package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"jcourse_go/dal"
	"jcourse_go/model/po"
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
