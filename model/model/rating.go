package model

type RatingRelatedType = string

const (
	RelatedTypeCourse       RatingRelatedType = "course"
	RelatedTypeTeacher      RatingRelatedType = "teacher"
	RelatedTypeTrainingPlan RatingRelatedType = "training_plan"
)

type RatingInfoDistItemByID struct {
	RelatedID int64
	Rating    int64
	Count     int64
}

type RatingInfoDistItem struct {
	Rating int64 `json:"rating"`
	Count  int64 `json:"count"`
}

type RatingInfo struct {
	MyRating   int64                `json:"my_rating"`
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
