package service

import (
	"context"

	"jcourse_go/dal"
	"jcourse_go/model/converter"
	"jcourse_go/model/dto"
	"jcourse_go/repository"
)

func CreateRating(ctx context.Context, userID int64, dto dto.RatingDTO) error {
	po := converter.ConvertRatingDTOToPO(userID, dto)
	query := repository.NewRatingQuery(dal.GetDBClient())
	err := query.CreateRating(ctx, po)
	if err != nil {
		return err
	}
	return nil
}
