package dto

import (
	model2 "jcourse_go/internal/model/model"
)

type TrainingPlanListResponse = BasePaginateResponse[model2.TrainingPlanSummary]

type TrainingPlanDetailResponse = model2.TrainingPlanDetail

type TrainingPlanDetailRequest struct {
	TrainingPlanID int64 `uri:"trainingPlanID" binding:"required"`
}
type TrainingPlanListQueryRequest struct {
	model2.PaginationFilterForQuery
	EntryYears  string `json:"entry_years" form:"entry_years"`
	Departments string `json:"departments" form:"departments"`
	Degrees     string `json:"degrees" form:"degrees"`
	MajorName   string `json:"major_name" form:"major_name"`
	MajorCode   string `json:"major_code" form:"major_code"`
}
