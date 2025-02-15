package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/dal"
	"jcourse_go/model/po"
	"jcourse_go/model/types"
)

func TestReviewQuery_GetReview(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewReviewQuery(db)

	reviews, err := query.GetReview(ctx, WithID(1))
	assert.Nil(t, err)
	assert.Len(t, reviews, 1)
	assert.Equal(t, int64(5), reviews[0].Rating)
}

func TestReviewQuery_GetReviewCount(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewReviewQuery(db)

	count, err := query.GetReviewCount(ctx)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), count)
}

func TestReviewQuery_CreateReview_normal(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewReviewQuery(db)
	ratingQuery := NewRatingQuery(db)

	t.Run("normal", func(t *testing.T) {
		reviewPO := po.ReviewPO{
			CourseID:    2,
			UserID:      2,
			Comment:     "",
			Rating:      5,
			Semester:    "",
			IsAnonymous: false,
		}
		id, err := query.CreateReview(ctx, reviewPO)
		assert.Nil(t, err)
		assert.NotZero(t, id)

		rating, err := ratingQuery.GetRatingInfo(ctx, types.RelatedTypeCourse, 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(2), rating.Count)
		assert.Equal(t, float64(5), rating.Average)
	})

}

func TestReviewQuery_UpdateReview(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewReviewQuery(db)
	ratingQuery := NewRatingQuery(db)

	t.Run("duplicate", func(t *testing.T) {
		reviewPO := po.ReviewPO{
			ID:          1,
			CourseID:    2,
			UserID:      1,
			Comment:     "",
			Rating:      1,
			Semester:    "",
			IsAnonymous: false,
		}

		err := query.UpdateReview(ctx, reviewPO)
		assert.NotNil(t, err)

		// no change
		reviews, err := query.GetReview(ctx, WithID(1))
		assert.Nil(t, err)
		assert.Equal(t, int64(5), reviews[0].Rating)

		info, err := ratingQuery.GetRatingInfo(ctx, types.RelatedTypeCourse, 1)
		if err != nil {
			return
		}
		assert.Len(t, info.RatingDist, 1)
		assert.Equal(t, int64(5), info.RatingDist[0].Rating)
	})

	t.Run("normal", func(t *testing.T) {
		reviewPO := po.ReviewPO{
			ID:          1,
			CourseID:    3,
			UserID:      1,
			Comment:     "",
			Rating:      1,
			Semester:    "",
			IsAnonymous: false,
		}

		err := query.UpdateReview(ctx, reviewPO)
		assert.Nil(t, err)
		reviews, err := query.GetReview(ctx, WithID(1))
		assert.Nil(t, err)
		assert.Equal(t, int64(1), reviews[0].Rating)

		info, err := ratingQuery.GetRatingInfo(ctx, types.RelatedTypeCourse, 3)
		if err != nil {
			return
		}
		assert.Equal(t, 1, len(info.RatingDist))
		if len(info.RatingDist) > 0 {
			assert.Equal(t, int64(1), info.RatingDist[0].Rating)
		}

		revision := po.ReviewRevisionPO{}
		err = db.Model(&po.ReviewRevisionPO{}).Where("review_id = ?", 1).Take(&revision).Error
		assert.Nil(t, err)
		assert.Equal(t, int64(5), revision.Rating)
	})

}

func TestReviewQuery_DeleteReview(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewReviewQuery(db)
	ratingQuery := NewRatingQuery(db)

	t.Run("normal", func(t *testing.T) {
		err := query.DeleteReview(ctx, WithID(1))
		assert.Nil(t, err)
		reviews, err := query.GetReview(ctx, WithID(1))
		assert.Nil(t, err)
		assert.Len(t, reviews, 0)

		info, err := ratingQuery.GetRatingInfo(ctx, types.RelatedTypeCourse, 1)
		if err != nil {
			return
		}
		assert.Len(t, info.RatingDist, 0)
	})
}
