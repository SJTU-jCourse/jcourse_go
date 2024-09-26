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

func GetTeacherDetail(ctx context.Context, teacherID int64) (*model.TeacherDetail, error) {
	if teacherID == 0 {
		return nil, errors.New("training-plan id is 0")
	}
	teacherQuery := repository.NewTeacherQuery(dal.GetDBClient())

	teacherPOs, err := teacherQuery.GetTeacher(ctx, repository.WithID(teacherID))
	if err != nil || len(teacherPOs) == 0 {
		return nil, err
	}
	teacherPO := teacherPOs[0]
	teacher := converter.ConvertTeacherDetailFromPO(teacherPO)

	courses, err := GetCourseList(ctx, model.CourseListFilterForQuery{MainTeacherID: teacherID})
	if err != nil {
		return nil, err
	}
	converter.PackTeacherWithCourses(&teacher, courses)

	ratingQuery := repository.NewRatingQuery(dal.GetDBClient())
	info, err := ratingQuery.GetRatingInfo(ctx, model.RelatedTypeTeacher, teacherID)
	if err != nil {
		return nil, err
	}
	converter.PackTeacherWithRatingInfo(&teacher.TeacherSummary, info)

	return &teacher, nil
}

func buildTeacherDBOptionFromFilter(query repository.ITeacherQuery, filter model.TeacherFilterForQuery) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.Name != "" {
		opts = append(opts, repository.WithName(filter.Name))
	}
	if filter.Code != "" {
		opts = append(opts, repository.WithCode(filter.Code))
	}
	if len(filter.Departments) > 0 {
		opts = append(opts, repository.WithDepartments(filter.Departments))
	}
	if len(filter.Titles) > 0 {
		opts = append(opts, repository.WithTitles(filter.Titles))
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
	if filter.PageSize > 0 {
		opts = append(opts, repository.WithLimit(filter.PageSize))
	}
	if filter.Page > 0 {
		opts = append(opts, repository.WithOffset(util.CalcOffset(filter.Page, filter.PageSize)))
	}
	return opts
}

func SearchTeacherList(ctx context.Context, filter model.TeacherFilterForQuery) ([]model.TeacherSummary, error) {
	teacherQuery := repository.NewTeacherQuery(dal.GetDBClient())
	t_opts := buildTeacherDBOptionFromFilter(teacherQuery, filter)

	teacherCourseQuery := repository.NewOfferedCourseQuery(dal.GetDBClient())
	validTeacherIDs, err := teacherCourseQuery.GetMainTeacherIDsWithOfferedCourseIDs(ctx, filter.ContainCourseIDs)
	if err != nil {
		return nil, err
	}
	t_opts = append(t_opts, repository.WithIDs(validTeacherIDs))

	teachers, err := teacherQuery.GetTeacher(ctx, t_opts...)
	if err != nil {
		return nil, err
	}

	teacherIDs := make([]int64, 0)
	for _, teacher := range teachers {
		teacherIDs = append(teacherIDs, int64(teacher.ID))
	}

	ratingQuery := repository.NewRatingQuery(dal.GetDBClient())
	infos, err := ratingQuery.GetRatingInfoByIDs(ctx, model.RelatedTypeTeacher, teacherIDs)
	if err != nil {
		return nil, err
	}

	domainTeachers := make([]model.TeacherSummary, 0)
	for _, t := range teachers {
		teacherDomain := converter.ConvertTeacherSummaryFromPO(t)
		converter.PackTeacherWithRatingInfo(&teacherDomain, infos[teacherDomain.ID])
		domainTeachers = append(domainTeachers, teacherDomain)
	}
	return domainTeachers, nil
}

func GetTeacherCount(ctx context.Context, filter model.TeacherFilterForQuery) (int64, error) {
	query := repository.NewTeacherQuery(dal.GetDBClient())
	filter.Page, filter.PageSize = 0, 0
	opts := buildTeacherDBOptionFromFilter(query, filter)
	return query.GetTeacherCount(ctx, opts...)
}

func GetTeacherFilter(ctx context.Context) (model.TeacherFilter, error) {
	query := repository.NewTeacherQuery(dal.GetDBClient())
	return query.GetTeacherFilter(ctx)
}
