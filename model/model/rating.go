package model

type RatingRelatedType = string

const (
	RelatedTypeCourse RatingRelatedType = "course"
)

type RatingInfoDistItemByID struct {
	RelatedID int64
	Rating    int64
	Count     int64
}

type RatingInfoDistItem struct {
	Rating int64
	Count  int64
}

type RatingInfo struct {
	Average    float64              `json:"average"`
	Count      int64                `json:"count"`
	RatingDist []RatingInfoDistItem `json:"rating_dist"`
}

func (r *RatingInfo) Calc() {
	if len(r.RatingDist) == 0 {
		return
	}
	for _, distItem := range r.RatingDist {
		r.Count = r.Count + distItem.Count
		r.Average = r.Average + float64(distItem.Count*distItem.Rating)
	}
	r.Average = r.Average / float64(r.Count)
}

type RatingDTO struct {
	RelatedType string `json:"related_type"`
	RelatedID   int64  `json:"related_id"`
}
