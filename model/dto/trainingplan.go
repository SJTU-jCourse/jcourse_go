package dto

import "jcourse_go/model/model"

type TrainingPlanListResponse = BasePaginateResponse[model.TrainingPlanSummary]

type TrainingPlanDetailResponse = model.TrainingPlanDetail

type TrainingPlanDetailRequest struct {
	TrainingPlanID int64 `uri:"trainingPlanID" binding:"required"`
}
type TrainingPlanListQueryRequest struct {
	model.PaginationFilterForQuery
	EntryYears  string `json:"entry_years" form:"entry_years"`
	Departments string `json:"departments" form:"departments"`
	Degrees     string `json:"degrees" form:"degrees"`
	MajorName   string `json:"major_name" form:"major_name"`
	MajorCode   string `json:"major_code" form:"major_code"`
}
