package service

import (
	"context"
	"errors"

	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/repository"
	"jcourse_go/util"
)

func buildReviewDBOptionFromFilter(ctx context.Context, q *repository.Query, filter model.ReviewFilterForQuery) repository.IReviewPODo {

	builder := q.ReviewPO.WithContext(ctx)
	r := q.ReviewPO

	if filter.Page > 0 || filter.PageSize > 0 {
		builder = builder.Offset(int(util.CalcOffset(filter.Page, filter.PageSize))).Limit(int(filter.PageSize))
	}
	if filter.Order != "" {
		field, ok := r.GetFieldByName(filter.Order)
		if ok {
			if filter.Ascending {
				builder = builder.Order(field)
			} else {
				builder = builder.Order(field.Desc())
			}
		}
	}

	if filter.CourseID != 0 {
		builder = builder.Where(r.CourseID.Eq(filter.CourseID))
	}
	if len(filter.Semester) > 0 {
		builder = builder.Where(r.Semester.Eq(filter.Semester))
	}
	if filter.UserID != 0 {
		builder = builder.Where(r.UserID.Eq(filter.UserID))
	}
	if filter.ReviewID != 0 {
		builder = builder.Where(r.ID.Eq(filter.ReviewID))
	}
	if filter.Rating != 0 {
		builder = builder.Where(r.Rating.Eq(filter.Rating))
	}
	if filter.ExcludeAnonymous {
		builder = builder.Where(r.IsAnonymous.Is(false))
	}
	return builder
}

func GetReviewList(ctx context.Context, currentUser *model.UserDetail, filter model.ReviewFilterForQuery) ([]model.Review, error) {
	r := repository.Q.ReviewPO
	q := buildReviewDBOptionFromFilter(ctx, repository.Q, filter)

	reviewPOs, err := q.Preload(r.Course, r.User, r.Reaction).Find()
	if err != nil {
		return nil, err
	}

	currentUserID := int64(0)
	if currentUser != nil {
		currentUserID = currentUser.ID
	}

	result := make([]model.Review, 0)
	for _, reviewPO := range reviewPOs {
		review := converter.ConvertReviewFromPO(*reviewPO)
		converter.PackReviewWithReaction(&review, currentUserID, reviewPO.Reaction)
		result = append(result, review)
	}

	return result, nil
}

func GetReviewCount(ctx context.Context, filter model.ReviewFilterForQuery) (int64, error) {
	filter.Page, filter.PageSize = 0, 0
	q := buildReviewDBOptionFromFilter(ctx, repository.Q, filter)
	return q.Count()
}

func CreateReview(ctx context.Context, review dto.UpdateReviewDTO, user *model.UserDetail) (int64, error) {
	if !validateReview(ctx, review, user) {
		return 0, errors.New("validate review error")
	}

	reviewPO := converter.ConvertReviewDTOToPO(review, user.ID)

	r := repository.Q.ReviewPO

	err := r.WithContext(ctx).Create(&reviewPO)
	if err != nil {
		return 0, err
	}
	return reviewPO.ID, nil
}

func UpdateReview(ctx context.Context, review dto.UpdateReviewDTO, user *model.UserDetail) error {
	if review.ID == 0 {
		return errors.New("no review id")
	}
	if !validateReview(ctx, review, user) {
		return errors.New("validate review error")
	}
	reviewPO := converter.ConvertReviewDTOToPO(review, user.ID)

	r := repository.Q.ReviewPO

	err := r.WithContext(ctx).Save(&reviewPO)
	if err != nil {
		return err
	}
	return nil
}

func DeleteReview(ctx context.Context, reviewID int64, user *model.UserDetail) error {
	r := repository.Q.ReviewPO

	review, err := r.WithContext(ctx).Where(r.ID.Eq(reviewID)).Take()
	if err != nil {
		return err
	}
	if user != nil && review.UserID != user.ID {
		return errors.New("no permission to delete review")
	}

	_, err = r.WithContext(ctx).Delete(review)
	if err != nil {
		return err
	}
	return nil
}

func validateReview(ctx context.Context, review dto.UpdateReviewDTO, user *model.UserDetail) bool {
	// 1. validate course and semester exists

	oc := repository.Q.OfferedCoursePO
	offeredCourse, err := oc.WithContext(ctx).Where(oc.CourseID.Eq(review.CourseID), oc.Semester.Eq(review.Semester)).Take()
	if err != nil || offeredCourse == nil {
		return false
	}

	// 2. validate comment

	// 3. validate review frequency

	return true
}
