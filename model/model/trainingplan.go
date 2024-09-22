package model

type TrainingPlanFilter struct {
	Page             int64
	PageSize         int64
	Major            string
	Department       string
	EntryYear        string
	ContainCourseIDs []int64
	SearchQuery      string
}

type TrainingPlanCourse struct {
	BaseCourse      BaseCourse `json:"base_course"`
	ID              int64      `json:"id"`
	SuggestSemester string     `json:"suggest_semester"`
	Category        string     `json:"category"`
}

type TrainingPlanSummary struct {
	ID         int64      `json:"id"`
	Code       string     `json:"code"`
	MajorName  string     `json:"name"`
	EntryYear  string     `json:"entry_year"`
	Department string     `json:"department"`
	Degree     string     `json:"degree"`
	RatingInfo RatingInfo `json:"rating_info"`
}

type TrainingPlanDetail struct {
	TrainingPlanSummary
	MajorClass string               `json:"major_class"`
	MinCredits float64              `json:"min_credits"`
	TotalYear  int64                `json:"total_year"`
	Courses    []TrainingPlanCourse `json:"courses"`
}
