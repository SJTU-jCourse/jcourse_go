package service

import (
	"context"
	"errors"

	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
	"jcourse_go/util"
)

func GetCourseDetail(ctx context.Context, courseID int64) (*domain.Course, error) {
	if courseID == 0 {
		return nil, errors.New("course id is 0")
	}
	courseQuery := repository.NewCourseQuery()
	coursePO, err := courseQuery.GetCourse(ctx, courseQuery.WithID(courseID))
	if err != nil {
		return nil, err
	}

	courseCategories, err := courseQuery.GetCourseCategories(ctx, []int64{int64(coursePO.ID)})
	if err != nil {
		return nil, err
	}

	teacherQuery := repository.NewTeacherQuery()
	teacherPO, err := teacherQuery.GetTeacher(ctx, teacherQuery.WithID(coursePO.MainTeacherID))
	if err != nil {
		return nil, err
	}

	offeredCourseQuery := repository.NewOfferedCourseQuery()
	offeredCoursePOs, err := offeredCourseQuery.GetOfferedCourseList(ctx, offeredCourseQuery.WithCourseID(courseID), offeredCourseQuery.WithOrderBy("semester", false))
	if err != nil {
		return nil, err
	}

	reviewQuery := repository.NewReviewQuery()
	infos, err := reviewQuery.GetCourseReviewInfo(ctx, []int64{courseID})
	if err != nil {
		return nil, err
	}

	course := converter.ConvertCoursePOToDomain(*coursePO)
	converter.PackCourseWithMainTeacher(&course, *teacherPO)
	converter.PackCourseWithOfferedCourse(&course, offeredCoursePOs)
	converter.PackCourseWithCategories(&course, courseCategories[course.ID])
	converter.PackCourseWithReviewInfo(&course, infos[course.ID])
	return &course, nil
}

func buildCourseDBOptionFromFilter(query repository.ICourseQuery, filter domain.CourseListFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.PageSize > 0 {
		opts = append(opts, query.WithLimit(filter.PageSize))
	}
	if filter.Page > 0 {
		opts = append(opts, query.WithOffset(util.CalcOffset(filter.Page, filter.PageSize)))
	}
	if len(filter.Categories) > 0 {
		opts = append(opts, query.WithCategories(filter.Categories))
	}
	if len(filter.Departments) > 0 {
		opts = append(opts, query.WithDepartments(filter.Departments))
	}
	if len(filter.Credits) > 0 {
		opts = append(opts, query.WithCredits(filter.Credits))
	}
	if filter.SearchQuery != "" {
		opts = append(opts, query.WithSearch(filter.SearchQuery))
	}
	return opts
}

func GetCourseList(ctx context.Context, filter domain.CourseListFilter) ([]domain.Course, error) {
	query := repository.NewCourseQuery()
	opts := buildCourseDBOptionFromFilter(query, filter)

	coursePOs, err := query.GetCourseList(ctx, opts...)
	if err != nil {
		return nil, err
	}

	courseIDs := make([]int64, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		courseIDs = append(courseIDs, int64(coursePO.ID))
	}

	courseCategories, err := query.GetCourseCategories(ctx, courseIDs)
	if err != nil {
		return nil, err
	}

	reviewQuery := repository.NewReviewQuery()
	infos, err := reviewQuery.GetCourseReviewInfo(ctx, courseIDs)
	if err != nil {
		return nil, err
	}

	courses := make([]domain.Course, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		course := converter.ConvertCoursePOToDomain(coursePO)
		converter.PackCourseWithCategories(&course, courseCategories[int64(coursePO.ID)])
		converter.PackCourseWithReviewInfo(&course, infos[int64(coursePO.ID)])
		courses = append(courses, course)
	}
	return courses, nil
}

func GetCourseCount(ctx context.Context, filter domain.CourseListFilter) (int64, error) {
	query := repository.NewCourseQuery()
	filter.Page, filter.PageSize = 0, 0
	opts := buildCourseDBOptionFromFilter(query, filter)
	return query.GetCourseCount(ctx, opts...)
}

func GetCourseByIDs(ctx context.Context, courseIDs []int64) (map[int64]domain.Course, error) {
	result := make(map[int64]domain.Course)
	if len(courseIDs) == 0 {
		return result, nil
	}
	courseQuery := repository.NewCourseQuery()
	courseMap, err := courseQuery.GetCourseByIDs(ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	for _, course := range courseMap {
		result[int64(course.ID)] = converter.ConvertCoursePOToDomain(course)
	}
	return result, nil
}
