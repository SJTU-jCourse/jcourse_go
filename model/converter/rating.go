package converter

import (
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func ConvertRatingInfoFromPO(po po.RatingPO) model.RatingInfo {
	return model.RatingInfo{}
}

func BuildRatingFromReview(review po.ReviewPO) po.RatingPO {
	return po.RatingPO{
		UserID:      review.UserID,
		RelatedType: model.RelatedTypeCourse,
		RelatedID:   review.CourseID,
		Rating:      review.Rating,
	}
}
