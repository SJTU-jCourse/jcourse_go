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

func GetTeacherDetail(ctx context.Context, teacherID int64) (*model.TeacherDetail, error) {
	if teacherID == 0 {
		return nil, errors.New("teacher id is 0")
	}
	t := repository.Q.TeacherPO
	teacherPO, err := t.WithContext(ctx).Where(t.ID.Eq(teacherID)).Take()
	if err != nil {
		return nil, err
	}
	if teacherPO == nil {
		return nil, errors.New("teacher not found")
	}
	teacher := converter.ConvertTeacherDetailFromPO(*teacherPO)

	courses, err := GetCourseList(ctx, model.CourseListFilterForQuery{MainTeacherID: teacherID})
	if err != nil {
		return nil, err
	}
	converter.PackTeacherWithCourses(&teacher, courses)

	info, err := GetRating(ctx, types.RelatedTypeTeacher, teacherID)
	if err != nil {
		return nil, err
	}
	converter.PackTeacherWithRatingInfo(&teacher.TeacherSummary, info)

	return &teacher, nil
}

func buildTeacherDBOptionFromFilter(ctx context.Context, q *repository.Query, filter model.TeacherFilterForQuery) repository.ITeacherPODo {
	builder := q.TeacherPO.WithContext(ctx)
	t := q.TeacherPO

	if filter.Page > 0 || filter.PageSize > 0 {
		builder = builder.Offset(int(util.CalcOffset(filter.Page, filter.PageSize))).Limit(int(filter.PageSize))
	}
	if filter.Order != "" {
		field, ok := t.GetFieldByName(filter.Order)
		if ok {
			if filter.Ascending {
				builder = builder.Order(field)
			} else {
				builder = builder.Order(field.Desc())
			}
		}
	}

	if filter.Name != "" {
		builder = builder.Where(t.Name.Eq(filter.Name))
	}
	if filter.Code != "" {
		builder = builder.Where(t.Code.Eq(filter.Code))
	}
	if len(filter.Departments) > 0 {
		builder = builder.Where(t.Department.In(filter.Departments...))
	}
	if len(filter.Titles) > 0 {
		builder = builder.Where(t.Title.In(filter.Titles...))
	}
	if filter.Pinyin != "" {
		builder = builder.Where(t.Pinyin.Eq(filter.Pinyin))
	}
	if filter.PinyinAbbr != "" {
		builder = builder.Where(t.PinyinAbbr.Eq(filter.PinyinAbbr))
	}
	return builder
}

func SearchTeacherList(ctx context.Context, filter model.TeacherFilterForQuery) ([]model.TeacherSummary, error) {
	q := buildTeacherDBOptionFromFilter(ctx, repository.Q, filter)

	teachers, err := q.Find()
	if err != nil {
		return nil, err
	}

	teacherIDs := make([]int64, 0)
	for _, teacher := range teachers {
		teacherIDs = append(teacherIDs, teacher.ID)
	}

	infos, err := GetMultipleRating(ctx, types.RelatedTypeTeacher, teacherIDs)
	if err != nil {
		return nil, err
	}

	domainTeachers := make([]model.TeacherSummary, 0)
	for _, t := range teachers {
		teacherDomain := converter.ConvertTeacherSummaryFromPO(*t)
		converter.PackTeacherWithRatingInfo(&teacherDomain, infos[teacherDomain.ID])
		domainTeachers = append(domainTeachers, teacherDomain)
	}
	return domainTeachers, nil
}

func GetTeacherCount(ctx context.Context, filter model.TeacherFilterForQuery) (int64, error) {

	filter.Page, filter.PageSize = 0, 0
	q := buildTeacherDBOptionFromFilter(ctx, repository.Q, filter)
	return q.Count()
}

func GetTeacherFilter(ctx context.Context) (model.TeacherFilter, error) {
	filter := model.TeacherFilter{
		Departments: make([]model.FilterItem, 0),
		Titles:      make([]model.FilterItem, 0),
	}

	t := repository.Q.TeacherPO
	err := t.WithContext(ctx).Group(t.Title.As("value"), t.ID.Count().As("count")).Scan(&filter.Titles)
	if err != nil {
		return filter, err
	}

	err = t.WithContext(ctx).Group(t.Department.As("value"), t.ID.Count().As("count")).Scan(&filter.Departments)
	if err != nil {
		return filter, err
	}
	return filter, nil
}
