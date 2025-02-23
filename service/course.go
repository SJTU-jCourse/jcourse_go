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

func GetCourseDetail(ctx context.Context, courseID int64, userID int64) (*model.CourseDetail, error) {
	if courseID == 0 {
		return nil, errors.New("course id is 0")
	}

	c := repository.Q.CoursePO
	coursePO, err := c.WithContext(ctx).Preload(c.MainTeacher, c.OfferedCourses, c.Categories).Where(c.ID.Eq(courseID)).Take()
	if err != nil {
		return nil, err
	}

	info, err := GetRating(ctx, types.RelatedTypeCourse, courseID)
	if err != nil {
		return nil, err
	}

	if userID != 0 {
		info.MyRating, _ = GetUserRating(ctx, types.RelatedTypeCourse, courseID, userID)
	}

	course := converter.ConvertCourseDetailFromPO(coursePO)
	converter.PackCourseWithRatingInfo(&course.CourseSummary, info)
	return &course, nil
}

func buildCourseDBOptionFromFilter(ctx context.Context, q *repository.Query, filter model.CourseListFilterForQuery) repository.ICoursePODo {
	builder := q.CoursePO.WithContext(ctx)
	c := q.CoursePO

	if filter.Page > 0 || filter.PageSize > 0 {
		builder = builder.Offset(int(util.CalcOffset(filter.Page, filter.PageSize))).Limit(int(filter.PageSize))
	}
	if filter.Order != "" {
		field, ok := q.CoursePO.GetFieldByName(filter.Order)
		if ok {
			if filter.Ascending {
				builder = builder.Order(field)
			} else {
				builder = builder.Order(field.Desc())
			}
		}
	}

	if len(filter.Categories) > 0 {
		builder = builder.Where(q.CourseCategoryPO.Category.In(filter.Categories...))
	}
	if len(filter.Departments) > 0 {
		builder = builder.Where(c.Department.In(filter.Departments...))
	}
	if len(filter.Credits) > 0 {
		builder = builder.Where(c.Credit.In(filter.Credits...))
	}
	if filter.Code != "" {
		builder = builder.Where(c.Code.Eq(filter.Code))
	}
	if filter.MainTeacherID > 0 {
		builder = builder.Where(c.MainTeacherID.Eq(filter.MainTeacherID))
	}

	return builder
}

func GetCourseList(ctx context.Context, filter model.CourseListFilterForQuery) ([]model.CourseSummary, error) {

	q := buildCourseDBOptionFromFilter(ctx, repository.Q, filter)

	coursePOs, err := q.Find()
	if err != nil {
		return nil, err
	}

	courseIDs := make([]int64, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		courseIDs = append(courseIDs, int64(coursePO.ID))
	}

	ratingMap, err := GetMultipleRating(ctx, types.RelatedTypeCourse, courseIDs)
	if err != nil {
		return nil, err
	}

	courses := make([]model.CourseSummary, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		course := converter.ConvertCourseSummaryFromPO(coursePO)
		converter.PackCourseWithRatingInfo(&course, ratingMap[coursePO.ID])
		courses = append(courses, course)
	}
	return courses, nil
}

func GetCourseCount(ctx context.Context, filter model.CourseListFilterForQuery) (int64, error) {
	filter.Page, filter.PageSize = 0, 0
	q := buildCourseDBOptionFromFilter(ctx, repository.Q, filter)
	return q.Count()
}

func GetBaseCourse(ctx context.Context, code string) (*model.BaseCourse, error) {
	c := repository.Q.BaseCoursePO
	baseCoursePO, err := c.WithContext(ctx).Where(c.Code.Eq(code)).Take()
	if err != nil {
		return nil, err
	}
	baseCourse := converter.ConvertBaseCourseFromPO(baseCoursePO)
	return &baseCourse, nil
}

func GetCourseFilter(ctx context.Context) (model.CourseFilter, error) {
	filter := model.CourseFilter{
		Categories:  make([]model.FilterItem, 0),
		Departments: make([]model.FilterItem, 0),
		Credits:     make([]model.FilterItem, 0),
		Semesters:   make([]model.FilterItem, 0),
	}

	c := repository.Q.CoursePO

	err := c.WithContext(ctx).Group(c.Credit.As("value"), c.ID.Count().As("count")).Scan(&filter.Credits)
	if err != nil {
		return filter, err
	}

	err = c.WithContext(ctx).Group(c.Department.As("value"), c.ID.Count().As("count")).Scan(&filter.Departments)
	if err != nil {
		return filter, err
	}

	oc := repository.Q.OfferedCoursePO
	err = oc.WithContext(ctx).Group(oc.Semester.As("value"), oc.CourseID.Count().As("count")).Scan(&filter.Semesters)
	if err != nil {
		return filter, err
	}

	cc := repository.Q.CourseCategoryPO
	err = cc.WithContext(ctx).Group(cc.Category.As("value"), cc.CourseID.Count().As("count")).Scan(&filter.Categories)
	if err != nil {
		return filter, err
	}
	return filter, nil
}
