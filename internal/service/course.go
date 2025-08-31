package service

import (
	"context"
	"errors"

	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/internal/model/converter"
	"jcourse_go/internal/model/types"
	"jcourse_go/pkg/util"
)

func GetRelatedCourses(ctx context.Context, course *course.CourseDetail) (*course.RelatedCourse, error) {
	if course == nil {
		return nil, errors.New("course is nil")
	}
	coursesUnderSameTeacher, err := GetCourseList(ctx, course.CourseListFilterForQuery{MainTeacherID: course.MainTeacher.ID})
	if err != nil {
		return nil, err
	}

	coursesWithOtherTeacher, err := GetCourseList(ctx, course.CourseListFilterForQuery{Code: course.Code})
	if err != nil {
		return nil, err
	}

	relatedCourse := &course.RelatedCourse{
		CoursesUnderSameTeacher:  make([]course.CourseSummary, 0),
		CoursesWithOtherTeachers: make([]course.CourseSummary, 0),
	}

	for _, c := range coursesUnderSameTeacher {
		if c.ID == course.ID {
			continue
		}
		relatedCourse.CoursesUnderSameTeacher = append(relatedCourse.CoursesUnderSameTeacher, c)
	}

	for _, c := range coursesWithOtherTeacher {
		if c.ID == course.ID {
			continue
		}
		relatedCourse.CoursesWithOtherTeachers = append(relatedCourse.CoursesWithOtherTeachers, c)
	}

	return relatedCourse, nil
}

func GetCourseDetail(ctx context.Context, courseID int64, userID int64) (*course.CourseDetail, error) {
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
	course.RelatedCourses, err = GetRelatedCourses(ctx, &course)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func buildCourseDBOptionFromFilter(ctx context.Context, q *repository.Query, filter course.CourseListQuery) repository.ICoursePODo {
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

func GetCourseList(ctx context.Context, filter course.CourseListQuery) ([]course.CourseSummary, error) {

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

	courses := make([]course.CourseSummary, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		course := converter.ConvertCourseSummaryFromPO(coursePO)
		converter.PackCourseWithRatingInfo(&course, ratingMap[coursePO.ID])
		courses = append(courses, course)
	}
	return courses, nil
}

func GetCourseCount(ctx context.Context, filter course.CourseListQuery) (int64, error) {
	filter.Page, filter.PageSize = 0, 0
	q := buildCourseDBOptionFromFilter(ctx, repository.Q, filter)
	return q.Count()
}

func GetBaseCourse(ctx context.Context, code string) (*course.BaseCourse, error) {
	c := repository.Q.BaseCoursePO
	baseCoursePO, err := c.WithContext(ctx).Where(c.Code.Eq(code)).Take()
	if err != nil {
		return nil, err
	}
	baseCourse := converter.ConvertBaseCourseFromPO(baseCoursePO)
	return &baseCourse, nil
}

func GetCourseFilter(ctx context.Context) (course.CourseFilter, error) {
	filter := course.CourseFilter{
		Categories:  make([]course.FilterAggregateItem, 0),
		Departments: make([]course.FilterAggregateItem, 0),
		Credits:     make([]course.FilterAggregateItem, 0),
		Semesters:   make([]course.FilterAggregateItem, 0),
	}

	c := repository.Q.CoursePO

	err := c.WithContext(ctx).Select(c.Credit.As("value"), c.ID.Count().As("count")).Group(c.Credit).Scan(&filter.Credits)
	if err != nil {
		return filter, err
	}

	err = c.WithContext(ctx).Select(c.Department.As("value"), c.ID.Count().As("count")).Group(c.Department).Scan(&filter.Departments)
	if err != nil {
		return filter, err
	}

	oc := repository.Q.OfferedCoursePO
	err = oc.WithContext(ctx).Select(oc.Semester.As("value"), oc.CourseID.Count().As("count")).Group(oc.Semester).Scan(&filter.Semesters)
	if err != nil {
		return filter, err
	}

	cc := repository.Q.CourseCategoryPO
	err = cc.WithContext(ctx).Select(cc.Category.As("value"), cc.CourseID.Count().As("count")).Group(cc.Category).Scan(&filter.Categories)
	if err != nil {
		return filter, err
	}
	return filter, nil
}
