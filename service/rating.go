package service

import (
	"context"

	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/model/types"
	"jcourse_go/query"
)

func CreateRating(ctx context.Context, userID int64, dto dto.RatingDTO) error {
	po := converter.ConvertRatingDTOToPO(userID, dto)
	r := query.Q.RatingPO
	err := r.WithContext(ctx).Create(&po)
	if err != nil {
		return err
	}
	return nil
}

func GetRating(ctx context.Context, relatedType types.RatingRelatedType, relatedID int64) (model.RatingInfo, error) {
	res := model.RatingInfo{}
	dist := make([]model.RatingInfoDistItem, 0)

	r := query.Q.RatingPO
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

func GetMultipleRating(ctx context.Context, relatedType types.RatingRelatedType, relatedIDs []int64) (map[int64]model.RatingInfo, error) {
	res := make(map[int64]model.RatingInfo)
	dist := make(map[int64][]model.RatingInfoDistItem)

	r := query.Q.RatingPO
	rows := make([]struct {
		RelatedID int64 `json:"related_id"`
		Rating    int64 `json:"rating"`
		Count     int64 `json:"count"`
	}, 0)
	err := r.WithContext(ctx).Select(r.RelatedID, r.Rating, r.ID.Count().As("count")).
		Where(r.RelatedID.In(relatedIDs...), r.RelatedType.Eq(string(relatedType))).
		Group(r.RelatedID, r.Rating).Scan(&rows)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		dist[row.RelatedID] = append(dist[row.RelatedID], model.RatingInfoDistItem{
			Rating: row.Rating,
			Count:  row.Count,
		})
	}

	for id, distItems := range dist {
		ratingInfo := model.RatingInfo{RatingDist: distItems}
		ratingInfo.Calc()
		res[id] = ratingInfo
	}

	return res, nil
}

func GetUserRating(ctx context.Context, relatedType types.RatingRelatedType, relatedID int64, userID int64) (int64, error) {
	r := query.Q.RatingPO
	rating, err := r.WithContext(ctx).Select(r.Rating).Where(r.RelatedID.Eq(relatedID), r.RelatedType.Eq(string(relatedType)), r.UserID.Eq(userID)).Take()
	if err != nil {
		return 0, err
	}
	return rating.Rating, nil
}
