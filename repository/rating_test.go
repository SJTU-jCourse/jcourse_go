package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"jcourse_go/dal"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func TestRatingQuery_GetRatingInfo(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewRatingQuery(db)
	info, err := query.GetRatingInfo(ctx, model.RelatedTypeCourse, 1)
	assert.Nil(t, err)
	assert.Len(t, info.RatingDist, 1)
	assert.Equal(t, int64(1), info.Count)
	assert.Equal(t, float64(5), info.Average)
}

func TestRatingQuery_GetRatingInfoByIDs(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewRatingQuery(db)
	infoMap, err := query.GetRatingInfoByIDs(ctx, model.RelatedTypeCourse, []int64{1})
	assert.Nil(t, err)
	assert.Len(t, infoMap, 1)
	info := infoMap[1]
	assert.Len(t, info.RatingDist, 1)
	assert.Equal(t, int64(1), info.Count)
	assert.Equal(t, float64(5), info.Average)
}

func TestRatingQuery_CreateRating(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewRatingQuery(db)
	rating := po.RatingPO{
		UserID:      3,
		RelatedType: model.RelatedTypeTeacher,
		RelatedID:   1,
		Rating:      5,
	}
	err := query.CreateRating(ctx, rating)
	assert.Nil(t, err)

	info, err := query.GetRatingInfo(ctx, model.RelatedTypeTeacher, 1)
	assert.Nil(t, err)
	assert.Len(t, info.RatingDist, 1)
	assert.Equal(t, int64(1), info.Count)
	assert.Equal(t, float64(5), info.Average)
}

func TestRatingQuery_UpdateRating(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewRatingQuery(db)
	rating := po.RatingPO{
		Model:       gorm.Model{ID: 3},
		UserID:      1,
		RelatedType: model.RelatedTypeTeacher,
		RelatedID:   2,
		Rating:      3,
	}
	err := query.UpdateRating(ctx, rating)
	assert.Nil(t, err)

	info, err := query.GetRatingInfo(ctx, model.RelatedTypeTeacher, 2)
	assert.Nil(t, err)
	assert.Len(t, info.RatingDist, 1)
	assert.Equal(t, int64(1), info.Count)
	assert.Equal(t, float64(3), info.Average)
}

func TestRatingQuery_DeleteRating(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewRatingQuery(db)
	rating := po.RatingPO{
		UserID:      1,
		RelatedType: model.RelatedTypeTeacher,
		RelatedID:   2,
		Rating:      3,
	}
	err := query.DeleteRating(ctx, rating)
	assert.Nil(t, err)

	info, err := query.GetRatingInfo(ctx, model.RelatedTypeTeacher, 2)
	assert.Nil(t, err)
	assert.Len(t, info.RatingDist, 0)
	assert.Equal(t, int64(0), info.Count)
}
