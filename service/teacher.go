package service

import (
	"context"
	"errors"
	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
)

func GetTeacherDetail(ctx context.Context, teacherID int64) (*domain.Teacher, error) {
	if teacherID == 0 {
		return nil, errors.New("training-plan id is 0")
	}
	teacherQuery := repository.NewTeacherQuery()

	teacherPO, err := teacherQuery.GetTeacher(ctx, teacherQuery.WithID(teacherID))
	if err != nil {
		return nil, err
	}
	teacher := converter.ConvertTeacherPOToDomain(teacherPO)

	courseQuery := repository.NewOfferedCourseQuery()
	courses, err := courseQuery.GetOfferedCourseList(ctx, courseQuery.WithMainTeacherID(teacherID))
	if err != nil {
		return nil, err
	}

	converter.PackTeacherWithCourses(teacher, courses)
	return teacher, nil
}

func buildTeacherDBOptionFromFilter(query repository.ITeacherQuery, filter domain.TeacherListFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.Name != "" {
		opts = append(opts, query.WithName(filter.Name))
	}
	if filter.Code != "" {
		opts = append(opts, query.WithCode(filter.Code))
	}
	if filter.Department != "" {
		opts = append(opts, query.WithDepartment(filter.Department))
	}
	if filter.Title != "" {
		opts = append(opts, query.WithTitle(filter.Title))
	}
	if filter.Pinyin != "" {
		opts = append(opts, query.WithPinyin(filter.Pinyin))
	}
	if filter.PinyinAbbr != "" {
		opts = append(opts, query.WithPinyinAbbr(filter.PinyinAbbr))
	}

	opts = append(opts, query.WithPaginate(filter.Page, filter.PageSize))
	return opts
}

func SearchTeacherList(ctx context.Context, filter domain.TeacherListFilter) ([]domain.Teacher, error) {
	teacherQuery := repository.NewTeacherQuery()
	t_opts := buildTeacherDBOptionFromFilter(teacherQuery, filter)

	teacherCourseQuery := repository.NewOfferedCourseQuery()
	validTeacherIDs, err := teacherCourseQuery.GetMainTeacherIDsWithOfferedCourseIDs(ctx, filter.ContainCourseIDs)
	if err != nil {
		return nil, err
	}
	t_opts = append(t_opts, teacherQuery.WithIDs(validTeacherIDs))

	teachers, err := teacherQuery.GetTeacherList(ctx, t_opts...)
	if err != nil {
		return nil, err
	}

	domainTeachers := make([]domain.Teacher, 0)
	for _, t := range teachers {
		q := repository.NewOfferedCourseQuery()
		offeredCoursePOs, err := q.GetOfferedCourseList(ctx, q.WithMainTeacherID(int64(t.ID)))
		if err != nil {
			return nil, err
		}
		teacherDomain := *converter.ConvertTeacherPOToDomain(&t)
		converter.PackTeacherWithCourses(&teacherDomain, offeredCoursePOs)
		domainTeachers = append(domainTeachers, teacherDomain)
	}
	return domainTeachers, nil
}

func GetTeacherCount(ctx context.Context, filter domain.TeacherListFilter) (int64, error) {
	query := repository.NewTeacherQuery()
	filter.Page, filter.PageSize = 0, 0
	opts := buildTeacherDBOptionFromFilter(query, filter)
	return query.GetTeacherCount(ctx, opts...)
}

func GetTeacherListByIDs(ctx context.Context, teacherIDs []int64) (map[int64]domain.Teacher, error) {

	teacherQuery := repository.NewTeacherQuery()
	teachers, err := teacherQuery.GetTeacherList(ctx, teacherQuery.WithIDs(teacherIDs))
	if err != nil {
		return nil, err
	}

	domainTeachers := make(map[int64]domain.Teacher)
	for _, t := range teachers {
		domainTeachers[int64(t.ID)] = *converter.ConvertTeacherPOToDomain(&t)
	}
	return domainTeachers, nil
}
