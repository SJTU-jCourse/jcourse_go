package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/model/po"
)

func ConvertReviewPOToDomain(review po.ReviewPO) domain.Review {
	return domain.Review{
		ID:          int64(review.ID),
		Course:      domain.Course{ID: review.CourseID},
		User:        domain.User{ID: review.UserID},
		Semester:    review.Semester,
		Rate:        review.Rate,
		IsAnonymous: review.IsAnonymous,
		Comment:     review.Comment,
		CreatedAt:   review.CreatedAt,
		UpdatedAt:   review.UpdatedAt,
	}
}

func PackReviewWithCourse(review *domain.Review, course domain.Course) {
	if review == nil {
		return
	}
	review.Course = course
}

func PackReviewWithUser(review *domain.Review, user domain.User) {
	if review == nil {
		return
	}
	review.User = user
}

func ConvertReviewDomainToDTO(review domain.Review, hideUser bool) dto.ReviewDTO {
	reviewDTO := dto.ReviewDTO{
		ID:          review.ID,
		Course:      ConvertCourseDomainToListDTO(review.Course),
		Semester:    review.Semester,
		Rate:        review.Rate,
		IsAnonymous: review.IsAnonymous,
		Comment:     review.Comment,
		UpdatedAt:   review.UpdatedAt,
		CreatedAt:   review.CreatedAt,
	}
	if !hideUser || !review.IsAnonymous {
		reviewDTO.User = ConvertUserDomainToReviewDTO(review.User)
	}
	return reviewDTO
}

func ConvertReviewDomainToListDTO(reviews []domain.Review, hideUser bool) []dto.ReviewDTO {
	result := make([]dto.ReviewDTO, 0)
	for _, review := range reviews {
		result = append(result, ConvertReviewDomainToDTO(review, hideUser))
	}
	return result
}

func ConvertUpdateReviewDTOToPO(review dto.UpdateReviewDTO, userID int64) po.ReviewPO {
	reviewPO := po.ReviewPO{
		CourseID:    review.CourseID,
		UserID:      userID,
		Comment:     review.Comment,
		Rate:        review.Rate,
		Semester:    review.Semester,
		IsAnonymous: review.IsAnonymous,
	}
	if review.ID != 0 {
		reviewPO.ID = uint(review.ID)
	}
	return reviewPO
}
