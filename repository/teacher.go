package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

type ITeacherQuery interface {
	GetTeacher(ctx context.Context, opts ...DBOption) ([]po.TeacherPO, error)
	GetTeacherCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetTeacherFilter(ctx context.Context) (model.TeacherFilter, error)
}

type TeacherQuery struct {
	db *gorm.DB
}

func (q *TeacherQuery) GetTeacherFilter(ctx context.Context) (model.TeacherFilter, error) {
	filter := model.TeacherFilter{
		Departments: make([]model.FilterItem, 0),
		Titles:      make([]model.FilterItem, 0),
	}
	err := q.db.WithContext(ctx).Model(&po.TeacherPO{}).
		Select("department as value, count(*) as count").
		Group("department").Find(&filter.Departments).Error
	if err != nil {
		return filter, err
	}
	err = q.db.WithContext(ctx).Model(&po.TeacherPO{}).
		Select("title as value, count(*) as count").
		Group("title").Find(&filter.Titles).Error
	if err != nil {
		return filter, err
	}
	return filter, nil
}

func NewTeacherQuery(db *gorm.DB) ITeacherQuery {
	return &TeacherQuery{
		db: db,
	}
}

func (q *TeacherQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := q.db.WithContext(ctx).Model(&po.TeacherPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (q *TeacherQuery) GetTeacher(ctx context.Context, opts ...DBOption) ([]po.TeacherPO, error) {
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
