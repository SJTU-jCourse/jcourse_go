package service

import (
	"context"

	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
	"jcourse_go/util"
)

func GetCourseList(ctx context.Context, filter domain.CourseListFilter) ([]domain.Course, error) {
	query := repository.NewCourseQuery()
	opts := make([]repository.DBOption, 0)

	opts = append(opts, query.WithLimit(filter.PageSize), query.WithOffset(util.CalcOffset(filter.Page, filter.PageSize)))
	if len(filter.Categories) > 0 {
		opts = append(opts, query.WithCategories(filter.Categories))
	}
	if len(filter.Departments) > 0 {
		opts = append(opts, query.WithDepartments(filter.Departments))
	}
	if len(filter.Credits) > 0 {
		opts = append(opts, query.WithCredits(filter.Credits))
	}

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

	courses := make([]domain.Course, 0, len(coursePOs))
	for _, coursePO := range coursePOs {
		courses = append(courses, converter.ConvertCoursePOToDomain(coursePO, courseCategories[int64(coursePO.ID)]))
	}
	return courses, nil
}

func GetCourseCount(ctx context.Context) (int64, error) {
	query := repository.NewCourseQuery()
	return query.GetCourseCount(ctx)
}
