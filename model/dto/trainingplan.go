package dto

type TrainingPlanCourseDTO struct {
	ID              int64   `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Credit          float64 `json:"credit"`
	SuggestYear     int64   `json:"suggest_year"`
	SuggestSemester int64   `json:"suggest_semester"`
	Department      string  `json:"department"`
}

type RateTrainingPlanRequest struct {
	TrainingPlanID int64  `form:"training_plan_id" binding:"required"`
	Rate           int64  `form:"rate" binding:"required"`
	Comment        string `form:"comment"`
}
type TrainingPlanReviewInfo struct {
	TrainingPlanID int64
	Avg            float64
	Count          int64
}
type TrainingPlanRateDTO struct {
	Rate  float64 `json:"rate"`
	Count int64   `json:"count"`
}
type TrainingPlanRateInfoDTO struct {
	Avg      float64               `json:"avg"`
	Count    int64                 `json:"count"`
	RateDist []TrainingPlanRateDTO `json:"rate_dist"`
}
type TrainingPlanListItemDTO struct {
	ID         int64                   `json:"id"`
	Code       string                  `json:"code"`
	MajorName  string                  `json:"name"`
	MinPoints  float64                 `json:"min_points"`
	MajorClass string                  `json:"major_class"`
	EntryYear  int64                   `json:"entry_year"`
	Department string                  `json:"department"`
	TotalYear  int64                   `json:"total_year"`
	Grade      float32                 `json:"grade"`
	Degree     string                  `json:"degree"`
	RateInfo   TrainingPlanRateInfoDTO `json:"rate_info"`
	Courses    []TrainingPlanCourseDTO `json:"courses"`
}

type TrainingPlanListResponse = BasePaginateResponse[TrainingPlanListItemDTO]
type TrainingPlanDetailResponse = TrainingPlanListItemDTO

type TrainingPlanDetailRequest struct {
	TrainingPlanID int64 `uri:"trainingPlanID" binding:"required"`
}
type TrainingPlanListQueryRequest struct {
	EntryYear     int64  `form:"entry_year"`
	Department    string `form:"department"`
	MajorName     string `form:"major_name"`
	MajorCode     string `form:"major_code"`
	SortDirection string `form:"sort_direction"`
	SortBy        string `form:"sort_by"`
	Page          int    `form:"page" binding:"required"`
	PageSize      int    `form:"page_size" binding:"required"`
}
type TrainingPlanListRequest struct {
	Page     int `form:"page" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}
