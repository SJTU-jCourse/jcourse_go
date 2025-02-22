package service

import (
	"context"
	"errors"

	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/model/types"
	"jcourse_go/repository"
	"jcourse_go/util"
)

func GetTrainingPlanDetail(ctx context.Context, trainingPlanID int64) (*model.TrainingPlanDetail, error) {
	if trainingPlanID == 0 {
		return nil, errors.New("training-plan id is 0")
	}

	tp := repository.Q.TrainingPlanPO
	trainingPlanPO, err := tp.WithContext(ctx).Preload(tp.BaseCourses, tp.BaseCourses.BaseCourse).Where(tp.ID.Eq(trainingPlanID)).Take()
	if err != nil {
		return nil, err
	}

	trainingPlan := converter.ConvertTrainingPlanDetailFromPO(*trainingPlanPO)

	info, err := GetRating(ctx, types.RelatedTypeTrainingPlan, trainingPlanID)
	if err != nil {
		return nil, err
	}
	converter.PackTrainingPlanWithRatingInfo(&trainingPlan.TrainingPlanSummary, info)

	return &trainingPlan, nil
}

func buildTrainingPlanDBOptionFromFilter(ctx context.Context, q *repository.Query, filter model.TrainingPlanFilterForQuery) repository.ITrainingPlanPODo {
	builder := q.TrainingPlanPO.WithContext(ctx)
	tp := q.TrainingPlanPO

	if filter.Page > 0 || filter.PageSize > 0 {
		builder = builder.Offset(int(util.CalcOffset(filter.Page, filter.PageSize))).Limit(int(filter.PageSize))
	}
	if filter.Order != "" {
		field, ok := tp.GetFieldByName(filter.Order)
		if ok {
			if filter.Ascending {
				builder = builder.Order(field)
			} else {
				builder = builder.Order(field.Desc())
			}
		}
	}

	if filter.Major != "" {
		builder = builder.Where(tp.Major.Eq(filter.Major))
	}
	if len(filter.EntryYears) > 0 {
		builder = builder.Where(tp.EntryYear.In(filter.EntryYears...))
	}
	if len(filter.Departments) > 0 {
		builder = builder.Where(tp.Department.In(filter.Departments...))
	}
	if len(filter.Degrees) > 0 {
		builder = builder.Where(tp.Degree.In(filter.Degrees...))
	}
	return builder
}

func GetTrainingPlanCount(ctx context.Context, filter model.TrainingPlanFilterForQuery) (int64, error) {
	filter.PageSize, filter.Page = 0, 0
	q := buildTrainingPlanDBOptionFromFilter(ctx, repository.Q, filter)
	return q.Count()
}

func SearchTrainingPlanList(ctx context.Context, filter model.TrainingPlanFilterForQuery) ([]model.TrainingPlanSummary, error) {

	q := buildTrainingPlanDBOptionFromFilter(ctx, repository.Q, filter)
	/*
		trainingPlanCourseQuery := repository.NewTrainingPlanCourseQuery(dal.GetDBClient())
		if len(filter.ContainCourseIDs) != 0 {
			tpc_opts := buildTrainingPlanCourseDBOptionFromFilter(trainingPlanCourseQuery, filter)
			validTrainingPlanIDs, err := trainingPlanCourseQuery.GetTrainingPlanListIDs(ctx, tpc_opts...)
			if err != nil {
				return nil, err
			}
			tp_opts = append(tp_opts, repository.WithIDs(validTrainingPlanIDs))
		}
	*/

	trainingPlanPOs, err := q.Find()
	if err != nil {
		return nil, err
	}

	trainingPlanIDs := make([]int64, 0)
	for _, tp := range trainingPlanPOs {
		trainingPlanIDs = append(trainingPlanIDs, tp.ID)
	}
	infos, err := GetMultipleRating(ctx, types.RelatedTypeTrainingPlan, trainingPlanIDs)
	if err != nil {
		return nil, err
	}
	result := make([]model.TrainingPlanSummary, 0)
	for _, tpPO := range trainingPlanPOs {
		tp := converter.ConvertTrainingPlanSummaryFromPO(*tpPO)
		converter.PackTrainingPlanWithRatingInfo(&tp, infos[tp.ID])
		result = append(result, tp)
	}
	return result, nil
}

func GetTrainingPlanFilter(ctx context.Context) (model.TrainingPlanFilter, error) {
	filter := model.TrainingPlanFilter{
		Departments: make([]model.FilterItem, 0),
		EntryYears:  make([]model.FilterItem, 0),
		Degrees:     make([]model.FilterItem, 0),
	}

	t := repository.Q.TrainingPlanPO
	err := t.WithContext(ctx).Group(t.Major.As("value"), t.ID.Count().As("count")).Scan(&filter.Degrees)
	if err != nil {
		return filter, err
	}

	err = t.WithContext(ctx).Group(t.Department.As("value"), t.ID.Count().As("count")).Scan(&filter.Departments)
	if err != nil {
		return filter, err
	}

	err = t.WithContext(ctx).Group(t.EntryYear.As("value"), t.ID.Count().As("count")).Scan(&filter.EntryYears)
	if err != nil {
		return filter, err
	}
	return filter, nil
}
