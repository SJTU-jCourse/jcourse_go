package service

import (
	"context"
	"errors"

	"jcourse_go/dal"
	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/repository"
	"jcourse_go/util"
)

func GetTrainingPlanDetail(ctx context.Context, trainingPlanID int64) (*model.TrainingPlanDetail, error) {
	if trainingPlanID == 0 {
		return nil, errors.New("training-plan id is 0")
	}
	trainingPlanQuery := repository.NewTrainingPlanQuery(dal.GetDBClient())

	trainingPlanPOs, err := trainingPlanQuery.GetTrainingPlan(ctx, repository.WithID(trainingPlanID))
	if err != nil || len(trainingPlanPOs) == 0 {
		return nil, err
	}
	trainingPlanPO := trainingPlanPOs[0]
	trainingPlan := converter.ConvertTrainingPlanDetailFromPO(trainingPlanPO)

	courseQuery := repository.NewTrainingPlanCourseQuery(dal.GetDBClient())
	courses, err := courseQuery.GetTrainingPlanCourseList(ctx, repository.WithTrainingPlanID(trainingPlanID))
	if err != nil {
		return nil, err
	}

	baseCourseIDs := make([]int64, 0)
	for _, course := range courses {
		baseCourseIDs = append(baseCourseIDs, course.BaseCourseID)
	}

	baseCourseQuery := repository.NewBaseCourseQuery(dal.GetDBClient())
	baseCoursePO, err := baseCourseQuery.GetBaseCoursesByIDs(ctx, baseCourseIDs)
	if err != nil {
		return nil, err
	}
	domainCourses := make([]model.TrainingPlanCourse, 0)
	for _, c := range courses {
		course := converter.ConvertTrainingPlanCourseFromPO(c)
		baseCourse := converter.ConvertBaseCourseFromPO(baseCoursePO[course.BaseCourse.ID])
		converter.PackTrainingPlanCourseWithBaseCourse(&course, baseCourse)
		domainCourses = append(domainCourses, course)
	}

	ratingQuery := repository.NewRatingQuery(dal.GetDBClient())
	info, err := ratingQuery.GetRatingInfo(ctx, model.RelatedTypeTrainingPlan, trainingPlanID)
	if err != nil {
		return nil, err
	}
	converter.PackTrainingPlanWithRatingInfo(&trainingPlan.TrainingPlanSummary, info)

	converter.PackTrainingPlanDetailWithCourse(&trainingPlan, domainCourses)
	return &trainingPlan, nil
}
func buildTrainingPlanDBOptionFromFilter(query repository.ITrainingPlanQuery, filter model.TrainingPlanFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.Major != "" {
		opts = append(opts, repository.WithMajor(filter.Major))
	}
	if filter.EntryYear != "" {
		opts = append(opts, repository.WithEntryYear(filter.EntryYear))
	}
	if filter.Department != "" {
		opts = append(opts, repository.WithDepartment(filter.Department))
	}
	if filter.SearchQuery != "" {
		opts = append(opts, repository.WithSearch(filter.SearchQuery))
	}
	if filter.PageSize > 0 {
		opts = append(opts, repository.WithLimit(filter.PageSize))
	}
	if filter.Page > 0 {
		opts = append(opts, repository.WithOffset(util.CalcOffset(filter.Page, filter.PageSize)))
	}
	return opts
}
func buildTrainingPlanCourseDBOptionFromFilter(query repository.ITrainingPlanCourseQuery, filter model.TrainingPlanFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if len(filter.ContainCourseIDs) > 0 {
		opts = append(opts, repository.WithCourseIDs(filter.ContainCourseIDs))
	}
	return opts
}
func GetTrainingPlanCount(ctx context.Context, filter model.TrainingPlanFilter) (int64, error) {
	trainingPlanQuery := repository.NewTrainingPlanQuery(dal.GetDBClient())
	filter.PageSize, filter.Page = 0, 0
	opts := buildTrainingPlanDBOptionFromFilter(trainingPlanQuery, filter)
	return trainingPlanQuery.GetTrainingPlanCount(ctx, opts...)
}

func SearchTrainingPlanList(ctx context.Context, filter model.TrainingPlanFilter) ([]model.TrainingPlanSummary, error) {

	trainingPlanQuery := repository.NewTrainingPlanQuery(dal.GetDBClient())
	tp_opts := buildTrainingPlanDBOptionFromFilter(trainingPlanQuery, filter)
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

	trainingPlanPOs, err := trainingPlanQuery.GetTrainingPlan(ctx, tp_opts...)
	if err != nil {
		return nil, err
	}

	trainingPlanIDs := make([]int64, 0)
	for _, tp := range trainingPlanPOs {
		trainingPlanIDs = append(trainingPlanIDs, int64(tp.ID))
	}

	ratingQuery := repository.NewRatingQuery(dal.GetDBClient())
	infos, err := ratingQuery.GetRatingInfoByIDs(ctx, model.RelatedTypeTrainingPlan, trainingPlanIDs)
	if err != nil {
		return nil, err
	}

	result := make([]model.TrainingPlanSummary, 0)
	for _, tpPO := range trainingPlanPOs {
		tp := converter.ConvertTrainingPlanSummaryFromPO(tpPO)
		converter.PackTrainingPlanWithRatingInfo(&tp, infos[tp.ID])
		result = append(result, tp)
	}
	return result, nil
}
