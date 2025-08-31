package converter

import (
	"jcourse_go/internal/domain/rating"
	entity2 "jcourse_go/internal/infrastructure/entity"
	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/model/types"
)

func ConvertRatingInfoFromPO(po entity2.RatingPO) rating.RatingInfo {
	return rating.RatingInfo{}
}

func BuildRatingFromReview(review entity2.ReviewPO) entity2.RatingPO {
	return entity2.RatingPO{
		UserID:      review.UserID,
		RelatedType: string(types.RelatedTypeCourse),
		RelatedID:   review.CourseID,
		Rating:      review.Rating,
	}
}

func ConvertRatingDTOToPO(userID int64, dto dto.RatingDTO) entity2.RatingPO {
	return entity2.RatingPO{
		UserID:      userID,
		RelatedType: dto.RelatedType,
		RelatedID:   dto.RelatedID,
		Rating:      dto.Rating,
	}
}
