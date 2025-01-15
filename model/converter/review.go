package converter

import (
	"gorm.io/gorm"

	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func ConvertReviewFromPO(po po.ReviewPO) model.Review {
	return model.Review{
		ID: int64(po.ID),
		Course: model.CourseMinimal{
			ID: po.CourseID,
		},
		User: model.UserMinimal{
			ID: po.UserID,
		},
		Comment:     po.Comment,
		Rating:      po.Rating,
		Semester:    po.Semester,
		IsAnonymous: po.IsAnonymous,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
		Grade:       po.Grade,
	}
}

func RemoveReviewUserInfo(review *model.Review, userID int64, hideUser bool) {
	if review == nil {
		return
	}
	// 本人点评不隐藏
	if hideUser && review.IsAnonymous && review.User.ID != userID {
		review.User = model.UserMinimal{}
	}
}

func RemoveReviewsUserInfo(reviews []model.Review, userID int64, hideUser bool) {
	for i := range reviews {
		RemoveReviewUserInfo(&reviews[i], userID, hideUser)
	}
}

func PackReviewWithCourse(review *model.Review, course model.CourseMinimal) {
	review.Course = course
}

func PackReviewWithUser(review *model.Review, user model.UserMinimal) {
	review.User = user
}

func ConvertReviewDTOToPO(dto dto.UpdateReviewDTO, userID int64) po.ReviewPO {
	return po.ReviewPO{
		Model: gorm.Model{
			ID: uint(dto.ID),
		},
		CourseID:    dto.CourseID,
		UserID:      userID,
		Comment:     dto.Comment,
		Rating:      dto.Rating,
		Semester:    dto.Semester,
		IsAnonymous: dto.IsAnonymous,
		Grade:       dto.Grade,
	}
}
