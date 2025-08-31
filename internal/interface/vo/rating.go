package vo

import "jcourse_go/internal/domain/rating"

type RatingVO struct {
	Average float64 `json:"average"`
	Count   int64   `json:"count"`
}

func NewRatingVO(r *rating.RatingInfo) RatingVO {
	return RatingVO{
		Average: r.Average,
		Count:   r.Count,
	}
}
