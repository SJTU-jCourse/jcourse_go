package service

import (
	"context"
	"errors"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
)

func GetTrainingPlanDetail(ctx context.Context, trainingPlanID int64) (*domain.TrainingPlanDetail, error) {
	if trainingPlanID == 0 {
		return nil, errors.New("training-plan id is 0")
	}
	trainingPlanQuery := repository.NewTrainingPlanQuery()

	trainingPlanPO, err := trainingPlanQuery.GetTrainingPlan(ctx, trainingPlanQuery.WithID(trainingPlanID))
	if err != nil {
		return nil, err
	}
	trainingPlan := converter.ConvertTrainingPlanPOToDomain(*trainingPlanPO)

	courseQuery := repository.NewTrainingPlanCourseQuery()
	courses, err := courseQuery.GetCourseListOfTrainingPlan(ctx, trainingPlanID)
	baseCourseQuery := repository.NewBaseCourseQuery()
	if err != nil {
		return nil, err
	}
	domainCourses := make([]domain.TrainingPlanCourse, 0)
	for _, c := range courses {
		baseCoursePO, err := baseCourseQuery.GetBaseCourse(ctx, baseCourseQuery.WithID(c.CourseID))
		if err != nil {
			return nil, err
		}
		courseDomain := converter.ConvertTrainingPlanCoursePOToDomain(c, *baseCoursePO)
		domainCourses = append(domainCourses, courseDomain)
	}
	converter.PackTrainingPlanDetailWithCourses(&trainingPlan, domainCourses)
	return &trainingPlan, nil
}
func buildTrainingPlanDBOptionFromFilter(query repository.ITrainingPlanQuery, filter domain.TrainingPlanFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.Major != "" {
		opts = append(opts, query.WithMajor(filter.Major))
	}
	if filter.EntryYear != "" {
		opts = append(opts, query.WithEntryYear(filter.EntryYear))
	}
	if filter.Department != "" {
		opts = append(opts, query.WithDepartment(filter.Department))
	}

	opts = append(opts, query.WithPaginate(filter.Page, filter.PageSize))
	return opts
}
func buildTrainingPlanCourseDBOptionFromFilter(query repository.ITrainingPlanCourseQuery, filter domain.TrainingPlanFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if len(filter.ContainCourseIDs) > 0 {
		opts = append(opts, query.WithCourseIDs(filter.ContainCourseIDs))
	}
	return opts
}
func GetTrainingPlanCount(ctx context.Context, filter domain.TrainingPlanFilter) int64 {
	trainingPlanQuery := repository.NewTrainingPlanQuery()
	opts := buildTrainingPlanDBOptionFromFilter(trainingPlanQuery, filter)
	return trainingPlanQuery.GetTrainingPlanCount(ctx, opts...)
}
func SearchTrainingPlanList(ctx context.Context, filter domain.TrainingPlanFilter) ([]domain.TrainingPlanDetail, error) {

	trainingPlanQuery := repository.NewTrainingPlanQuery()
	tp_opts := buildTrainingPlanDBOptionFromFilter(trainingPlanQuery, filter)
	trainingPlanCourseQuery := repository.NewTrainingPlanCourseQuery()
	if len(filter.ContainCourseIDs) != 0 {
		tpc_opts := buildTrainingPlanCourseDBOptionFromFilter(trainingPlanCourseQuery, filter)
		validTrainingPlanIDs, err := trainingPlanCourseQuery.GetTrainingPlanListIDs(ctx, tpc_opts...)
		if err != nil {
			return nil, err
		}
		tp_opts = append(tp_opts, trainingPlanQuery.WithIDs(validTrainingPlanIDs))
	}

	trainingPlanIDs, err := trainingPlanQuery.GetTrainingPlanListIDs(ctx, tp_opts...)
	if err != nil {
		return nil, err
	}
	result := make([]domain.TrainingPlanDetail, 0)
	for _, id := range trainingPlanIDs {
		d, err := GetTrainingPlanDetail(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, *d)
	}
	return result, nil
}

func GetTrainingPlanListByIDs(ctx context.Context, trainingPlanIDs []int64) (map[int64]domain.TrainingPlanDetail, error) {

	domainTrainingPlans := make(map[int64]domain.TrainingPlanDetail, 0)
	for _, id := range trainingPlanIDs {
		data, err := GetTrainingPlanDetail(ctx, id)
		if err != nil {
			return nil, err
		}
		domainTrainingPlans[id] = *data
	}
	return domainTrainingPlans, nil
}
