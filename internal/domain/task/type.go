package task

import "time"

type Type string

const (
	TypeReviewAntiSpam Type = "review.anti_spam"
	TypeReviewReward   Type = "review.reward"
)

type ReviewAntiSpamPayload struct {
	ReviewID  int64     `json:"review_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
