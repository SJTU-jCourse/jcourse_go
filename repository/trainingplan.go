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
	optionDB(ctx context.Context, opts ...DBOption) *gorm.DB
	GetTrainingPlan(ctx context.Context, opts ...DBOption) (*po.TrainingPlanPO, error)
	GetTrainingPlanList(ctx context.Context, opts ...DBOption) ([]po.TrainingPlanPO, error)
	GetTrainingPlanListIDs(ctx context.Context, opts ...DBOption) ([]int64, error)
	GetTrainingPlanCount(ctx context.Context, opts ...DBOption) int64
}

func NewTrainingPlanQuery(db *gorm.DB) ITrainingPlanQuery {
	return &TrainingPlanQuery{db: db}
}

func (t *TrainingPlanQuery) GetTrainingPlanListIDs(ctx context.Context, opts ...DBOption) ([]int64, error) {
	db := t.optionDB(ctx, opts...)
	var ids []int64
	result := db.Select("id").Find(&ids)
	if result.Error != nil {
		return nil, result.Error
	}
	return ids, nil
}
func (t *TrainingPlanQuery) GetTrainingPlanCount(ctx context.Context, opts ...DBOption) int64 {
	db := t.optionDB(ctx, opts...)
	var count int64
	result := db.Count(&count)
	if result.Error != nil {
		return 0
	}
	return count

}

func (t *TrainingPlanQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := t.db.Model(&po.TrainingPlanPO{}).WithContext(ctx)
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (t *TrainingPlanQuery) GetTrainingPlan(ctx context.Context, opts ...DBOption) (*po.TrainingPlanPO, error) {
	db := t.optionDB(ctx, opts...)
	trainingPlan := po.TrainingPlanPO{}
	result := db.Debug().First(&trainingPlan)
	if result.Error != nil {
		return nil, result.Error
	}
	return &trainingPlan, nil
}
func (t *TrainingPlanQuery) GetTrainingPlanList(ctx context.Context, opts ...DBOption) ([]po.TrainingPlanPO, error) {
	db := t.optionDB(ctx, opts...)
	trainingPlans := make([]po.TrainingPlanPO, 0)
	result := db.Debug().Find(&trainingPlans)
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
	GetCourseListOfTrainingPlan(ctx context.Context, trainingPlanID int64) ([]po.TrainingPlanCoursePO, error)
	GetTrainingPlanListIDs(ctx context.Context, opts ...DBOption) ([]int64, error)
	optionDB(ctx context.Context, opts ...DBOption) *gorm.DB
}

func (t *TrainingPlanCourseQuery) GetTrainingPlanListIDs(ctx context.Context, opts ...DBOption) ([]int64, error) {
	validTrainingPlans := make([]int64, 0)
	db := t.optionDB(ctx, opts...)
	result := db.Distinct("training_plan_id").Pluck("training_plan_id", &validTrainingPlans)
	if result.Error != nil {
		return nil, result.Error
	}
	return validTrainingPlans, nil
}
func (t *TrainingPlanCourseQuery) GetTrainingPlanCourseList(ctx context.Context, opts ...DBOption) ([]po.TrainingPlanCoursePO, error) {
	db := t.optionDB(ctx, opts...)
	trainingPlanCourses := make([]po.TrainingPlanCoursePO, 0)
	result := db.Debug().Find(&trainingPlanCourses)
	if result.Error != nil {
		return nil, result.Error
	}
	return trainingPlanCourses, nil
}

func (t *TrainingPlanCourseQuery) GetCourseListOfTrainingPlan(ctx context.Context, trainingPlanID int64) ([]po.TrainingPlanCoursePO, error) {
	db := t.optionDB(ctx, WithTrainingPlanID(trainingPlanID))
	courses := make([]po.TrainingPlanCoursePO, 0)
	result := db.Debug().Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}
	return courses, nil
}
