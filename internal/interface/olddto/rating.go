package olddto

type RatingDTO struct {
	RelatedType string `json:"related_type" binding:"required"`
	RelatedID   int64  `json:"related_id" binding:"required"`
	Rating      int64  `json:"rating" binding:"required"`
}
