package service

import (
	"context"
	"errors"

	"jcourse_go/model/converter"
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
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
	if filter.ReviewID != 0 {
		opts = append(opts, query.WithID(filter.ReviewID))
	}
	return opts
}

func GetReviewList(ctx context.Context, filter domain.ReviewFilter) ([]domain.Review, error) {
	reviewQuery := repository.NewReviewQuery()
	opts := buildReviewDBOptionFromFilter(reviewQuery, filter)
	reviewPOs, err := reviewQuery.GetReviewList(ctx, opts...)
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
		result = append(result, review)
	}

	return result, nil
}

func GetReviewCount(ctx context.Context, filter domain.ReviewFilter) (int64, error) {
	query := repository.NewReviewQuery()
	filter.Page, filter.PageSize = 0, 0
	opts := buildReviewDBOptionFromFilter(query, filter)
	return query.GetReviewCount(ctx, opts...)
}

func CreateReview(ctx context.Context, review dto.UpdateReviewDTO, user *domain.User) (int64, error) {
	if !validateReview(ctx, review, user) {
		return 0, errors.New("validate review error")
	}
	query := repository.NewReviewQuery()
	reviewPO := converter.ConvertUpdateReviewDTOToPO(review, user.ID)
	reviewID, err := query.CreateReview(ctx, reviewPO)
	if err != nil {
		return 0, err
	}
	return reviewID, nil
}

func UpdateReview(ctx context.Context, review dto.UpdateReviewDTO, user *domain.User) error {
	if review.ID == 0 {
		return errors.New("no review id")
	}
	if !validateReview(ctx, review, user) {
		return errors.New("validate review error")
	}
	query := repository.NewReviewQuery()
	reviewPO := converter.ConvertUpdateReviewDTOToPO(review, user.ID)
	_, err := query.UpdateReview(ctx, reviewPO)
	if err != nil {
		return err
	}
	return nil
}

func DeleteReview(ctx context.Context, reviewID int64) error {
	query := repository.NewReviewQuery()
	_, err := query.DeleteReview(ctx, query.WithID(reviewID))
	if err != nil {
		return err
	}
	return nil
}

func validateReview(ctx context.Context, review dto.UpdateReviewDTO, user *domain.User) bool {
	// 1. validate course and semester exists
	offeredCourseQuery := repository.NewOfferedCourseQuery()
	offeredCourse, err := offeredCourseQuery.GetOfferedCourse(ctx, offeredCourseQuery.WithCourseID(review.CourseID), offeredCourseQuery.WithSemester(review.Semester))
	if err != nil || offeredCourse == nil {
		return false
	}

	// 2. validate comment

	// 3. validate review frequency

	return true
}
