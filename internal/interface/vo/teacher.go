package vo

import (
	"jcourse_go/internal/domain/rating"
)

type TeacherSummary struct {
	ID         int64             `json:"id"`
	Name       string            `json:"name"`
	Department string            `json:"department"`
	Picture    string            `json:"picture"`
	RatingInfo rating.RatingInfo `json:"rating_info"`
}
