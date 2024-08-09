package dto

type RatingDTO struct {
	RelatedType string `json:"related_type"`
	RelatedID   int64  `json:"related_id"`
}
