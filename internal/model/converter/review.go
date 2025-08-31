package converter

import (
	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/user"
	"jcourse_go/internal/infrastructure/entity"
	"jcourse_go/internal/interface/dto"
)

func ConvertReviewFromPO(po *entity.ReviewPO) course.Review {
	review := course.Review{
		ID:          po.ID,
		Course:      ConvertCourseMinimalFromPO(po.Course),
		User:        ConvertUserMinimalFromPO(po.User),
		Comment:     po.Comment,
		Rating:      po.Rating,
		Semester:    po.Semester,
		IsAnonymous: po.IsAnonymous,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
		Score:       po.Grade,
	}
	return review
}

func RemoveReviewUserInfo(review *course.Review, userID int64, hideUser bool) {
	if review == nil {
		return
	}
	// 本人点评不隐藏
	if hideUser && review.IsAnonymous && review.User.ID != userID {
		review.User = user.UserMinimal{}
	}
}

func RemoveReviewsUserInfo(reviews []course.Review, userID int64, hideUser bool) {
	for i := range reviews {
		RemoveReviewUserInfo(&reviews[i], userID, hideUser)
	}
}

func ConvertReviewDTOToPO(dto dto.UpdateReviewDTO, userID int64) entity.ReviewPO {
	return entity.ReviewPO{
		ID:          dto.ID,
		CourseID:    dto.CourseID,
		UserID:      userID,
		Comment:     dto.Comment,
		Rating:      dto.Rating,
		Semester:    dto.Semester,
		IsAnonymous: dto.IsAnonymous,
		Grade:       dto.Grade,
	}
}

func PackReviewWithReaction(review *course.Review, currentUserID int64, reactions []entity.ReviewReactionPO) {
	if review == nil {
		return
	}
	if review.Reaction.TotalReactions == nil {
		review.Reaction.TotalReactions = make([]review.ReactionItem, 0)
	}
	if review.Reaction.MyReactions == nil {
		review.Reaction.MyReactions = make(map[string]int64)
	}

	reactionMap := make(map[string]int64)

	for _, reaction := range reactions {
		if reaction.UserID == currentUserID {
			review.Reaction.MyReactions[reaction.Reaction] = reaction.ID
		}
		reactionMap[reaction.Reaction] += 1
	}

	for reaction, count := range reactionMap {
		review.Reaction.TotalReactions = append(review.Reaction.TotalReactions, review.ReactionItem{
			Reaction: reaction,
			Count:    count,
		})
	}
}
