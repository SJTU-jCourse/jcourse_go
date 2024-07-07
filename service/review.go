package service

import (
	"context"

	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/repository"
	"jcourse_go/util"
)

func buildReviewDBOptionFromFilter(query repository.IReviewQuery, filter domain.ReviewFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.PageSize > 0 {
		opts = append(opts, query.WithLimit(filter.PageSize))
	}
	if filter.Page > 0 {
		opts = append(opts, query.WithOffset(util.CalcOffset(filter.Page, filter.PageSize)))
	}
	if filter.CourseID != 0 {
		opts = append(opts, query.WithCourseID(filter.CourseID))
	}
	if len(filter.Semester) > 0 {
		opts = append(opts, query.WithSemester(filter.Semester))
	}
	if filter.UserID != 0 {
		opts = append(opts, query.WithUserID(filter.UserID))
	}
	return opts
}

func GetReviewList(ctx context.Context, filter domain.ReviewFilter) ([]domain.Review, error) {
	reviewQuery := repository.NewCourseReviewQuery()
	opts := buildReviewDBOptionFromFilter(reviewQuery, filter)
	reviewPOs, err := reviewQuery.GetCourseReviewList(ctx, opts...)
	if err != nil {
		return nil, err
	}

	courseIDs := make([]int64, 0)
	userIDs := make([]int64, 0)

	for _, review := range reviewPOs {
		courseIDs = append(courseIDs, review.CourseID)
		userIDs = append(userIDs, review.UserID)
	}

	courseMap, err := GetCourseByIDs(ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	userMap, err := GetUserByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Review, 0)

	for _, reviewPO := range reviewPOs {
		review := converter.ConvertReviewPOToDomain(reviewPO)

		course, ok := courseMap[reviewPO.CourseID]
		if ok {
			converter.PackReviewWithCourse(&review, course)
		}
		user, ok := userMap[reviewPO.UserID]
		if ok {
			converter.PackReviewWithUser(&review, user)
		}
	}

	return result, nil
}

func GetReviewCount(ctx context.Context, filter domain.ReviewFilter) (int64, error) {
	query := repository.NewCourseReviewQuery()
	filter.Page, filter.PageSize = 0, 0
	opts := buildReviewDBOptionFromFilter(query, filter)
	return query.GetCourseReviewCount(ctx, opts...)
}

func CreateReview(ctx context.Context, review domain.Review) error {
	return nil
}
