package repository

import (
	"context"

	"jcourse_go/model/po"

	"gorm.io/gorm"
)

type TrainingPlanQuery struct {
	db *gorm.DB
}

type ITrainingPlanQuery interface {
	GetTrainingPlan(ctx context.Context, opts ...DBOption) ([]po.TrainingPlanPO, error)
	GetTrainingPlanCount(ctx context.Context, opts ...DBOption) (int64, error)
}

func NewTrainingPlanQuery(db *gorm.DB) ITrainingPlanQuery {
	return &TrainingPlanQuery{db: db}
}

func (t *TrainingPlanQuery) GetTrainingPlanCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := t.optionDB(ctx, opts...)
	var count int64
	result := db.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil

}

func (t *TrainingPlanQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := t.db.Model(&po.TrainingPlanPO{}).WithContext(ctx)
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (t *TrainingPlanQuery) GetTrainingPlan(ctx context.Context, opts ...DBOption) ([]po.TrainingPlanPO, error) {
	db := t.optionDB(ctx, opts...)
	trainingPlans := make([]po.TrainingPlanPO, 0)
	result := db.Find(&trainingPlans)
	if result.Error != nil {
		return nil, result.Error
	}
	return trainingPlans, nil
}

type TrainingPlanCourseQuery struct {
	db *gorm.DB
}

func NewTrainingPlanCourseQuery(db *gorm.DB) ITrainingPlanCourseQuery {
	return &TrainingPlanCourseQuery{db: db}
}

func (t *TrainingPlanCourseQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := t.db.Model(&po.TrainingPlanCoursePO{}).WithContext(ctx)
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

type ITrainingPlanCourseQuery interface {
	GetTrainingPlanCourseList(ctx context.Context, opts ...DBOption) ([]po.TrainingPlanCoursePO, error)
}

func (t *TrainingPlanCourseQuery) GetTrainingPlanCourseList(ctx context.Context, opts ...DBOption) ([]po.TrainingPlanCoursePO, error) {
	db := t.optionDB(ctx, opts...)
	trainingPlanCourses := make([]po.TrainingPlanCoursePO, 0)
	result := db.Find(&trainingPlanCourses)
	if result.Error != nil {
		return nil, result.Error
	}
	return trainingPlanCourses, nil
}
