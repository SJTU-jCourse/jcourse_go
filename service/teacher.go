package service

import (
	"context"
	"errors"

	"jcourse_go/dal"
	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/repository"
)

func GetTeacherDetail(ctx context.Context, teacherID int64) (*model.TeacherDetail, error) {
	if teacherID == 0 {
		return nil, errors.New("training-plan id is 0")
	}
	teacherQuery := repository.NewTeacherQuery(dal.GetDBClient())

	teacherPO, err := teacherQuery.GetTeacher(ctx, repository.WithID(teacherID))
	if err != nil {
		return nil, err
	}
	teacher := converter.ConvertTeacherDetailFromPO(*teacherPO)

	courseQuery := repository.NewOfferedCourseQuery(dal.GetDBClient())
	coursePOs, err := courseQuery.GetOfferedCourseList(ctx, repository.WithMainTeacherID(teacherID))
	if err != nil {
		return nil, err
	}
	courses := converter.ConvertOfferedCoursesFromPOs(coursePOs)
	converter.PackTeacherWithOfferedCourses(&teacher, courses)
	return &teacher, nil
}

func buildTeacherDBOptionFromFilter(query repository.ITeacherQuery, filter model.TeacherListFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.Name != "" {
		opts = append(opts, repository.WithName(filter.Name))
	}
	if filter.Code != "" {
		opts = append(opts, repository.WithCode(filter.Code))
	}
	if filter.Department != "" {
		opts = append(opts, repository.WithDepartment(filter.Department))
	}
	if filter.Title != "" {
		opts = append(opts, repository.WithTitle(filter.Title))
	}
	if filter.Pinyin != "" {
		opts = append(opts, repository.WithPinyin(filter.Pinyin))
	}
	if filter.PinyinAbbr != "" {
		opts = append(opts, repository.WithPinyinAbbr(filter.PinyinAbbr))
	}
	if filter.SearchQuery != "" {
		opts = append(opts, repository.WithSearch(filter.SearchQuery))
	}

	opts = append(opts, repository.WithPaginate(filter.Page, filter.PageSize))
	return opts
}

func SearchTeacherList(ctx context.Context, filter model.TeacherListFilter) ([]model.TeacherSummary, error) {
	teacherQuery := repository.NewTeacherQuery(dal.GetDBClient())
	t_opts := buildTeacherDBOptionFromFilter(teacherQuery, filter)

	teacherCourseQuery := repository.NewOfferedCourseQuery(dal.GetDBClient())
	validTeacherIDs, err := teacherCourseQuery.GetMainTeacherIDsWithOfferedCourseIDs(ctx, filter.ContainCourseIDs)
	if err != nil {
		return nil, err
	}
	t_opts = append(t_opts, repository.WithIDs(validTeacherIDs))

	teachers, err := teacherQuery.GetTeacherList(ctx, t_opts...)
	if err != nil {
		return nil, err
	}

	domainTeachers := make([]model.TeacherSummary, 0)
	for _, t := range teachers {
		teacherDomain := converter.ConvertTeacherSummaryFromPO(t)
		domainTeachers = append(domainTeachers, teacherDomain)
	}
	return domainTeachers, nil
}

func GetTeacherCount(ctx context.Context, filter model.TeacherListFilter) (int64, error) {
	query := repository.NewTeacherQuery(dal.GetDBClient())
	filter.Page, filter.PageSize = 0, 0
	opts := buildTeacherDBOptionFromFilter(query, filter)
	return query.GetTeacherCount(ctx, opts...)
}

func GetTeacherListByIDs(ctx context.Context, teacherIDs []int64) (map[int64]model.TeacherSummary, error) {
	teacherQuery := repository.NewTeacherQuery(dal.GetDBClient())
	teachers, err := teacherQuery.GetTeacherList(ctx, repository.WithIDs(teacherIDs))
	if err != nil {
		return nil, err
	}

	domainTeachers := make(map[int64]model.TeacherSummary)
	for _, t := range teachers {
		domainTeachers[int64(t.ID)] = converter.ConvertTeacherSummaryFromPO(t)
	}
	return domainTeachers, nil
}
