package converter

import (
	"jcourse_go/internal/model/dto"
	model2 "jcourse_go/internal/model/model"
	"jcourse_go/internal/model/po"
)

func ConvertReviewFromPO(po *po.ReviewPO) model2.Review {
	review := model2.Review{
		ID:          po.ID,
		Course:      ConvertCourseMinimalFromPO(po.Course),
		User:        ConvertUserMinimalFromPO(po.User),
		Comment:     po.Comment,
		Rating:      po.Rating,
		Semester:    po.Semester,
		IsAnonymous: po.IsAnonymous,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
		Grade:       po.Grade,
	}
	return review
}

func RemoveReviewUserInfo(review *model2.Review, userID int64, hideUser bool) {
	if review == nil {
		return
	}
	// 本人点评不隐藏
	if hideUser && review.IsAnonymous && review.User.ID != userID {
		review.User = model2.UserMinimal{}
	}
}

func RemoveReviewsUserInfo(reviews []model2.Review, userID int64, hideUser bool) {
	for i := range reviews {
		RemoveReviewUserInfo(&reviews[i], userID, hideUser)
	}
}

func ConvertReviewDTOToPO(dto dto.UpdateReviewDTO, userID int64) po.ReviewPO {
	return po.ReviewPO{
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

func PackReviewWithReaction(review *model2.Review, currentUserID int64, reactions []po.ReviewReactionPO) {
	if review == nil {
		return
	}
	if review.Reaction.TotalReactions == nil {
		review.Reaction.TotalReactions = make([]model2.ReactionItem, 0)
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
		review.Reaction.TotalReactions = append(review.Reaction.TotalReactions, model2.ReactionItem{
			Reaction: reaction,
			Count:    count,
		})
	}
}
