package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"jcourse_go/model/po"
)

// TODO: 暂时没有添加从ProfileDesc中搜索的功能
type ITeacherQuery interface {
	GetTeacher(ctx context.Context, opts ...DBOption) (*po.TeacherPO, error)
	GetTeacherList(ctx context.Context, opts ...DBOption) ([]po.TeacherPO, error)
	GetTeacherCount(ctx context.Context, opts ...DBOption) (int64, error)
}

type TeacherQuery struct {
	db *gorm.DB
}

func NewTeacherQuery(db *gorm.DB) ITeacherQuery {
	return &TeacherQuery{
		db: db,
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

func (q *TeacherQuery) GetTeacherCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := q.optionDB(ctx, opts...)
	var count int64
	result := db.Model(&po.TeacherPO{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
