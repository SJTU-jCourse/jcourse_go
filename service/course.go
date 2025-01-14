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

func GetCourseDetail(ctx context.Context, courseID int64) (*model.CourseDetail, error) {
	if courseID == 0 {
		return nil, errors.New("course id is 0")
	}
	courseQuery := repository.NewCourseQuery(dal.GetDBClient())
	coursePOs, err := courseQuery.GetCourse(ctx, repository.WithID(courseID))
	if err != nil || len(coursePOs) == 0 {
		return nil, err
	}
	coursePO := coursePOs[0]
	courseCategories, err := courseQuery.GetCourseCategories(ctx, []int64{int64(coursePO.ID)})
	if err != nil {
		return nil, err
	}

	teacherQuery := repository.NewTeacherQuery(dal.GetDBClient())
	teacherPOs, err := teacherQuery.GetTeacher(ctx, repository.WithID(coursePO.MainTeacherID))
	if err != nil || len(teacherPOs) == 0 {
		return nil, err
	}
	teacherPO := teacherPOs[0]

	offeredCourseQuery := repository.NewOfferedCourseQuery(dal.GetDBClient())
	offeredCoursePOs, err := offeredCourseQuery.GetOfferedCourse(ctx, repository.WithCourseID(courseID), repository.WithOrderBy("semester", false))
	if err != nil {
		return nil, err
	}

	ratingQuery := repository.NewRatingQuery(dal.GetDBClient())
	info, err := ratingQuery.GetRatingInfo(ctx, model.RelatedTypeCourse, courseID)
	if err != nil {
		return nil, err
	}

	course := converter.ConvertCourseDetailFromPO(coursePO)
	converter.PackCourseWithMainTeacher(&course.CourseMinimal, converter.ConvertTeacherSummaryFromPO(teacherPO))
	offeredCourses := converter.ConvertOfferedCoursesFromPOs(offeredCoursePOs)
	converter.PackCourseWithOfferedCourse(&course, offeredCourses)
	converter.PackCourseWithCategories(&course.CourseSummary, courseCategories[course.ID])
	converter.PackCourseWithRatingInfo(&course.CourseSummary, info)
	return &course, nil
}

func buildCourseDBOptionFromFilter(query repository.ICourseQuery, filter model.CourseListFilterForQuery) []repository.DBOption {
	opts := buildPaginationDBOptions(filter.PaginationFilterForQuery)
	if len(filter.Categories) > 0 {
		opts = append(opts, repository.WithCategories(filter.Categories))
	}
	if len(filter.Departments) > 0 {
		opts = append(opts, repository.WithDepartments(filter.Departments))
	}
	if len(filter.Credits) > 0 {
		opts = append(opts, repository.WithCredits(filter.Credits))
	}
	if filter.Code != "" {
		opts = append(opts, repository.WithCode(filter.Code))
	}
	if filter.MainTeacherID > 0 {
		opts = append(opts, repository.WithMainTeacherID(filter.MainTeacherID))
	}
	return opts
}

func GetCourseList(ctx context.Context, filter model.CourseListFilterForQuery) ([]model.CourseSummary, error) {
	query := repository.NewCourseQuery(dal.GetDBClient())
	opts := buildCourseDBOptionFromFilter(query, filter)

	coursePOs, err := query.GetCourse(ctx, opts...)
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

	ratingQuery := repository.NewRatingQuery(dal.GetDBClient())
	infos, err := ratingQuery.GetRatingInfoByIDs(ctx, model.RelatedTypeCourse, courseIDs)
	if err != nil {
		return nil, err
	}

	courses := make([]model.CourseSummary, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		course := converter.ConvertCourseSummaryFromPO(coursePO)
		converter.PackCourseWithCategories(&course, courseCategories[int64(coursePO.ID)])
		converter.PackCourseWithRatingInfo(&course, infos[int64(coursePO.ID)])
		courses = append(courses, course)
	}
	return courses, nil
}

func GetCourseCount(ctx context.Context, filter model.CourseListFilterForQuery) (int64, error) {
	query := repository.NewCourseQuery(dal.GetDBClient())
	filter.Page, filter.PageSize = 0, 0
	opts := buildCourseDBOptionFromFilter(query, filter)
	return query.GetCourseCount(ctx, opts...)
}

func GetCourseByIDs(ctx context.Context, courseIDs []int64) (map[int64]model.CourseSummary, error) {
	result := make(map[int64]model.CourseSummary)
	if len(courseIDs) == 0 {
		return result, nil
	}
	courseQuery := repository.NewCourseQuery(dal.GetDBClient())
	courseMap, err := courseQuery.GetCourseByIDs(ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	for _, course := range courseMap {
		result[int64(course.ID)] = converter.ConvertCourseSummaryFromPO(course)
	}
	return result, nil
}

func GetBaseCourse(ctx context.Context, code string) (*model.BaseCourse, error) {
	query := repository.NewBaseCourseQuery(dal.GetDBClient())
	baseCourses, err := query.GetBaseCourse(ctx, repository.WithCode(code))
	if err != nil {
		return nil, err
	}
	if len(baseCourses) == 0 {
		return nil, errors.New("no base course")
	}
	baseCourse := converter.ConvertBaseCourseFromPO(baseCourses[0])
	return &baseCourse, nil
}

func GetCourseFilter(ctx context.Context) (model.CourseFilter, error) {
	query := repository.NewCourseQuery(dal.GetDBClient())
	return query.GetCourseFilter(ctx)
}

func buildPaginationDBOptions(filter model.PaginationFilterForQuery) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.PageSize > 0 {
		opts = append(opts, repository.WithLimit(filter.PageSize))
	}
	if filter.Page > 0 {
		opts = append(opts, repository.WithOffset(util.CalcOffset(filter.Page, filter.PageSize)))
	}
	if filter.Search != "" {
		opts = append(opts, repository.WithSearch(filter.Search))
	}
	if filter.Order != "" {
		opts = append(opts, repository.WithOrderBy(filter.Order, filter.Ascending))
	}
	return opts
}
