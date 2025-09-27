package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type TrainingPlanQueryService interface {
	GetTrainingPlanList(ctx context.Context, query course.TrainingPlanListQuery) ([]vo.TrainingPlanVO, error)
	GetTrainingPlanDetail(ctx context.Context, trainingPlanID shared.IDType) (*vo.TrainingPlanDetailVO, error)
	GetTrainingPlanFilter(ctx context.Context) (*course.TrainingPlanFilter, error)
}

type trainingPlanQueryService struct {
	db *gorm.DB
}

func (t *trainingPlanQueryService) GetTrainingPlanList(ctx context.Context, query course.TrainingPlanListQuery) ([]vo.TrainingPlanVO, error) {
	tps := make([]entity.TrainingPlan, 0)
	db := t.db.WithContext(ctx).Model(&entity.TrainingPlan{})
	if query.Major != "" {
		db = db.Where("major = ?", query.Major)
	}
	if len(query.EntryYears) > 0 {
		db = db.Where("entry_year IN (?)", query.EntryYears)
	}
	if len(query.Degrees) > 0 {
		db = db.Where("degree IN (?)", query.Degrees)
	}
	if len(query.Departments) > 0 {
		db = db.Where("department IN (?)", query.Departments)
	}
	if err := db.Find(&tps).Error; err != nil {
		return nil, err
	}

	tpVOs := make([]vo.TrainingPlanVO, 0)
	for _, tp := range tps {
		tpVOs = append(tpVOs, vo.NewTrainingPlanVOFromEntity(&tp))
	}
	return tpVOs, nil
}

func (t *trainingPlanQueryService) GetTrainingPlanDetail(ctx context.Context, trainingPlanID shared.IDType) (*vo.TrainingPlanDetailVO, error) {
	tp := entity.TrainingPlan{}
	if err := t.db.WithContext(ctx).
		Model(&entity.TrainingPlan{}).
		Preload("Curriculums").
		Where("id = ?", trainingPlanID).
		Take(&tp).Error; err != nil {
		return nil, err
	}
	tpVO := vo.NewTrainingPlanDetailVOFromEntity(&tp)
	return &tpVO, nil
}

func (t *trainingPlanQueryService) GetTrainingPlanFilter(ctx context.Context) (*course.TrainingPlanFilter, error) {
	// TODO implement me
	panic("implement me")
}

func NewTrainingPlanQueryService(db *gorm.DB) TrainingPlanQueryService {
	return &trainingPlanQueryService{db: db}
}
