package converter

import (
	"jcourse_go/internal/domain/rating"
	entity2 "jcourse_go/internal/infrastructure/entity"
	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/model/types"
)

func ConvertRatingInfoFromPO(po entity2.Rating) rating.RatingInfo {
	return rating.RatingInfo{}
}

func BuildRatingFromReview(review entity2.Review) entity2.Rating {
	return entity2.Rating{
		UserID:      review.UserID,
		RelatedType: string(types.RelatedTypeCourse),
		RelatedID:   review.CourseID,
		Rating:      review.Rating,
	}
}

func ConvertRatingDTOToPO(userID int64, dto olddto.RatingDTO) entity2.Rating {
	return entity2.Rating{
		UserID:      userID,
		RelatedType: dto.RelatedType,
		RelatedID:   dto.RelatedID,
		Rating:      dto.Rating,
	}
}
