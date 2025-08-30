package converter

import (
	"jcourse_go/internal/model/dto"
	"jcourse_go/internal/model/model"
	po2 "jcourse_go/internal/model/po"
	"jcourse_go/internal/model/types"
)

func ConvertRatingInfoFromPO(po po2.RatingPO) model.RatingInfo {
	return model.RatingInfo{}
}

func BuildRatingFromReview(review po2.ReviewPO) po2.RatingPO {
	return po2.RatingPO{
		UserID:      review.UserID,
		RelatedType: string(types.RelatedTypeCourse),
		RelatedID:   review.CourseID,
		Rating:      review.Rating,
	}
}

func ConvertRatingDTOToPO(userID int64, dto dto.RatingDTO) po2.RatingPO {
	return po2.RatingPO{
		UserID:      userID,
		RelatedType: dto.RelatedType,
		RelatedID:   dto.RelatedID,
		Rating:      dto.Rating,
	}
}
