package dto

type TrainingPlanCourseDTO struct {
	ID              int64   `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Credit          float64 `json:"credit"`
	SuggestSemester string  `json:"suggest_semester"`
	Department      string  `json:"department"`
}

type TrainingPlanListItemDTO struct {
	ID         int64                   `json:"id"`
	Code       string                  `json:"code"`
	MajorName  string                  `json:"name"`
	MinCredits float64                 `json:"min_credits"`
	MajorClass string                  `json:"major_class"`
	EntryYear  int64                   `json:"entry_year"`
	Department string                  `json:"department"`
	TotalYear  int64                   `json:"total_year"`
	Degree     string                  `json:"degree"`
	Courses    []TrainingPlanCourseDTO `json:"courses"`
}

type TrainingPlanListResponse = BasePaginateResponse[TrainingPlanListItemDTO]
type TrainingPlanDetailResponse = TrainingPlanListItemDTO

type TrainingPlanDetailRequest struct {
	TrainingPlanID int64 `uri:"trainingPlanID" binding:"required"`
}
type TrainingPlanListQueryRequest struct {
	EntryYear     int64  `json:"entry_year" form:"entry_year"`
	Department    string `json:"department" form:"department"`
	MajorName     string `json:"major_name" form:"major_name"`
	MajorCode     string `json:"major_code" form:"major_code"`
	SortDirection string `json:"sort_direction" form:"sort_direction"`
	SortBy        string `json:"sort_by" form:"sort_by"`
	Page          int    `json:"page" binding:"required" form:"page"`
	PageSize      int    `json:"page_size" binding:"required" form:"page_size"`
	SearchQuery   string `json:"search_query" form:"search_query"`
}
type TrainingPlanListRequest struct {
	Page     int `json:"page" binding:"required" form:"page"`
	PageSize int `json:"page_size" binding:"required" form:"page_size"`
}
