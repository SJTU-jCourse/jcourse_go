package service

import (
	"context"

	"jcourse_go/dal"
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/model/types"
	"jcourse_go/query"
)

func CreateRating(ctx context.Context, userID int64, dto dto.RatingDTO) error {
	po := converter.ConvertRatingDTOToPO(userID, dto)
	r := query.Use(dal.GetDBClient()).RatingPO
	err := r.WithContext(ctx).Create(&po)
	if err != nil {
		return err
	}
	return nil
}

func GetRating(ctx context.Context, relatedType types.RatingRelatedType, relatedID int64) (model.RatingInfo, error) {
	res := model.RatingInfo{}
	dist := make([]model.RatingInfoDistItem, 0)

	r := query.Use(dal.GetDBClient()).RatingPO
	err := r.WithContext(ctx).Select(r.Rating, r.ID.Count().As("count")).
		Where(r.RelatedID.Eq(relatedID), r.RelatedType.Eq(string(relatedType))).
		Group(r.Rating).Scan(&dist)
	if err != nil {
		return res, err
	}
	res.RatingDist = dist
	res.Calc()
	return res, nil
}
