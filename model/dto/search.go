package dto

// shared between search types
type SearchRequest struct {
	Query    string          `json:"query" binding:"required"`
	Page     int64           `json:"page" binding:"required"`
	PageSize int64           `json:"page_size" binding:"required"`
	Filter   SearchFilterDTO `json:"filter"`
	SortBy   string          `json:"sort_by"`
}

type SearchFilterDTO struct {
	// should guarentee to exist in frontend
	Department   string         `json:"department"`
	TrainingPlan string         `json:"training_plan"`
	TeacherName  string         `json:"teacher"`
	ActiveYear   int            `json:"year"`
	RatingRange  RatingRangeDTO `json:"rating_range"`
}

type RatingRangeDTO struct {
	From int `json:"from" binding:"required"`
	To   int `json:"to" binding:"required"`
}

// response: use ReviewListResponse, CourseListResponse, and TrainingPlanListResponse
