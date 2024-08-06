package repository

import (
	"context"
	"jcourse_go/dal"
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
	WithID(id int64) DBOption
	WithDepartment(department string) DBOption
	WithMajor(major string) DBOption
	WithEntryYear(entryYear string) DBOption
	WithDegree(degree string) DBOption
	WithIDs(courseIDs []int64) DBOption
	WithPaginate(page int64, pageSize int64) DBOption
}

func NewTrainingPlanQuery() ITrainingPlanQuery {
	return &TrainingPlanQuery{db: dal.GetDBClient()}
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
func (t *TrainingPlanQuery) WithPaginate(page int64, pageSize int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 || pageSize <= 0 {
			return db.Where("1 = 0")
		}
		return db.Offset(int((page - 1) * pageSize)).Limit(int(pageSize))
	}
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

func (t *TrainingPlanQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (t *TrainingPlanQuery) WithDepartment(department string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("department = ?", department)
	}
}

func (t *TrainingPlanQuery) WithMajor(major string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("major = ?", major)
	}
}

func (t *TrainingPlanQuery) WithEntryYear(entryYear string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("entry_year = ?", entryYear)
	}
}

func (t *TrainingPlanQuery) WithDegree(degree string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("degree = ?", degree)
	}
}

type TrainingPlanCourseQuery struct {
	db *gorm.DB
}

func NewTrainingPlanCourseQuery() ITrainingPlanCourseQuery {
	return &TrainingPlanCourseQuery{db: dal.GetDBClient()}
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
	WithTrainingPlanID(trainingPlanID int64) DBOption
	WithCourseID(courseID int64) DBOption
	WithCourseIDs(courseIDs []int64) DBOption
	WithSuggestSemester(semester int64) DBOption
	WithSuggestYear(year int64) DBOption
	WithDepartment(department string) DBOption
}

func (t *TrainingPlanCourseQuery) WithSuggestYear(year int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("suggest_year = ?", year)
	}
}
func (t *TrainingPlanCourseQuery) WithSuggestSemester(semester int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("suggest_semester = ?", semester)
	}
}
func (t *TrainingPlanCourseQuery) WithDepartment(department string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("department = ?", department)
	}
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

func (t *TrainingPlanCourseQuery) WithCourseID(courseID int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id = ?", courseID)
	}
}
func (t *TrainingPlanCourseQuery) WithTrainingPlanID(trainingPlanID int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("training_plan_id = ?", trainingPlanID)
	}
}
func (t *TrainingPlanCourseQuery) WithCourseIDs(courseIDs []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id IN ?", courseIDs).
			Group("training_plan_id").
			Having("count(DISTINCT course_id) = ?", len(courseIDs))
	}
}
func (t *TrainingPlanCourseQuery) GetCourseListOfTrainingPlan(ctx context.Context, trainingPlanID int64) ([]po.TrainingPlanCoursePO, error) {
	db := t.optionDB(ctx, t.WithTrainingPlanID(trainingPlanID))
	courses := make([]po.TrainingPlanCoursePO, 0)
	result := db.Debug().Find(&courses)
	if result.Error != nil {
		return nil, result.Error
	}
	return courses, nil
}
func (t *TrainingPlanQuery) WithIDs(trainingPlanIDs []int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id IN ?", trainingPlanIDs)
	}
}

type ITrainingPlanReviewQuery interface {
	GetTrainingPlanRateInfo(ctx context.Context, trainingPlanIDs []int64) ([]po.TrainingPlanRatePO, error)
	GetTrainingPlanReviewList(ctx context.Context, ops ...DBOption) ([]po.TrainingPlanPO, error)
	optionDB(ctx context.Context, opts ...DBOption) *gorm.DB
	WithUserID(id int64) DBOption
	WithTrainingPlanID(id int64) DBOption
}
type TrainingPlanReviewQuery struct {
	db *gorm.DB
}

func (t *TrainingPlanReviewQuery) WithTrainingPlanID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("training_plan_id = ?", id)
	}
}
func (t *TrainingPlanReviewQuery) WithUserID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", id)
	}
}
func (t *TrainingPlanReviewQuery) OptionDB(ctx context.Context, ops ...DBOption) *gorm.DB {
	db := t.db.Model(&po.TrainingPlanReviewPO{}).WithContext(ctx)
	for _, opt := range ops {
		db = opt(db)
	}
	return db
}
func (t *TrainingPlanReviewQuery) GetTrainingPlanRateInfo(ctx context.Context, trainingPlanIDs []int64) ([]po.TrainingPlanRatePO, error) {
	db := t.OptionDB(ctx)
	db.Group("training_plan_id").Select("count(*) as count, avg(rate) as average, course_id")
	rates := make([]po.TrainingPlanRatePO, 0)

	result := db.Debug().Where("training_plan_id IN ?", trainingPlanIDs).Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}
	return rates, nil
}

func (t *TrainingPlanReviewQuery) GetTrainingPlanReviewList(ctx context.Context, ops ...DBOption) ([]po.TrainingPlanPO, error) {
	db := t.OptionDB(ctx, ops...)
	trainingPlans := make([]po.TrainingPlanPO, 0)
	result := db.Debug().Find(&trainingPlans)
	if result.Error != nil {
		return nil, result.Error
	}
	return trainingPlans, nil
}
