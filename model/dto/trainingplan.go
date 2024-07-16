package dto

type TrainingPlanCourseDTO = BaseCourseDTO


type TrainingPlanListItemDTO struct {
	ID          int64   `json:"id"`
	Code        string  `json:"code"`
	MajorName 	string 	`json:"name"`
	MinPoints 	float32 `json:"min_points"`
	MajorClass  string `json:"major_class"`
	EntryYear   int64  `json:"entry_year"`
	Department  string `json:"department"`
	TotalYear   int64  `json:"total_year"`
	Grade       float32 `json:"grade"`
	Courses     []TrainingPlanCourseDTO `json:"courses"`
}

type TrainingPlanListResponse = BasePaginateResponse[TrainingPlanListItemDTO]
type TrainingPlanDetailResponse = TrainingPlanListItemDTO

type TrainingPlanDetailRequest struct {
	TrainingPlanID int64 `uri:"trainingPlanID" binding:"required"`
}
type TrainingPlanListQueryRequest struct {
	EntryYear int64 `json:"entry_year"`
	Department string `json:"department"`
	MajorName string `json:"major_name"`
	MajorCode string `json:"major_code"`
	SortDirection string `json:"sort_direction"`
	SortBy string `json:"sort_by"`
	Page  int `json:"page"`
	PageSize int `json:"page_size"`
}