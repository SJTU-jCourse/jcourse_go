package service

import (
	"context"
	"errors"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
)

func GetTrainingPlanDetail(ctx context.Context, trainingPlanID int64) (*domain.TrainingPlan, error) {
	if trainingPlanID == 0{
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
	if err != nil {
		return nil, err
	}
	domainBaseCourses := make([]domain.BaseCourse, 0)
	for _,c := range courses{
		domainBaseCourses = append(domainBaseCourses, converter.ConvertBaseCoursePOToDomain(c))
	}
	converter.PackTrainingPlanWithCourses(&trainingPlan, domainBaseCourses)
	return &trainingPlan, nil
}
func buildTrainingPlanDBOptionFromFilter(query repository.ITrainingPlanQuery,filter domain.TrainingPlanFilter) []repository.DBOption{
	opts := make([]repository.DBOption, 0)
	if filter.Major != "" {
		opts = append(opts, query.WithMajor(filter.Major))
	}
	if filter.EntryYear != ""{
		opts = append(opts, query.WithEntryYear(filter.EntryYear))
	}
	if filter.Department != ""{
		opts = append(opts, query.WithDepartment(filter.Department))
	}
	return opts
}
func buildTrainingPlanCourseDBOptionFromFilter(query repository.ITrainingPlanCourseQuery, filter domain.TrainingPlanFilter) []repository.DBOption{
	opts := make([]repository.DBOption, 0)
	if len(filter.ContainCourseIDs) > 0{
		opts = append(opts, query.WithCourseIDs(filter.ContainCourseIDs))
	}
	return opts
}
func SearchTrainingPlanList(ctx context.Context, filter domain.TrainingPlanFilter) ([]domain.TrainingPlan, error){

	trainingPlanQuery := repository.NewTrainingPlanQuery()
	tp_opts := buildTrainingPlanDBOptionFromFilter(trainingPlanQuery, filter)

	trainingPlanCourseQuery := repository.NewTrainingPlanCourseQuery()
	tpc_opts := buildTrainingPlanCourseDBOptionFromFilter(trainingPlanCourseQuery, filter)
	validTrainingPlanIDs, err := trainingPlanCourseQuery.GetTrainingPlanListIDs(ctx, tpc_opts...)
	if err != nil {
		return nil, err
	}
	tp_opts = append(tp_opts, trainingPlanQuery.WithIDs(validTrainingPlanIDs))

	trainingPlans, err := trainingPlanQuery.GetTrainingPlanList(ctx, tp_opts...)
	if err != nil {
		return nil,err
	}

	domainTrainingPlans := make([]domain.TrainingPlan, 0)
	for _, tp := range trainingPlans{
		domainTrainingPlans = append(domainTrainingPlans, converter.ConvertTrainingPlanPOToDomain(tp))
	}
	return domainTrainingPlans, nil
}

func GetTrainingPlanListByIDs(ctx context.Context, trainingPlanIDs []int64) (map[int64]domain.TrainingPlan, error){

	trainingPlanQuery := repository.NewTrainingPlanQuery()
	trainingPlans, err := trainingPlanQuery.GetTrainingPlanList(ctx, trainingPlanQuery.WithIDs(trainingPlanIDs))
	if err != nil {
		return nil, err
	}
	domainTrainingPlans := make(map[int64]domain.TrainingPlan)
	for _, tp := range trainingPlans{
		domainTrainingPlans[int64(tp.ID)] = converter.ConvertTrainingPlanPOToDomain(tp)
	}
	return domainTrainingPlans, nil
}