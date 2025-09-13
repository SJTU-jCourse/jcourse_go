package olddto

import (
	"jcourse_go/internal/domain/course"
)

type TrainingPlanListResponse = BasePaginateResponse[course.TrainingPlanSummary]

type TrainingPlanDetailResponse = course.TrainingPlanDetail

type TrainingPlanDetailRequest struct {
	TrainingPlanID int64 `uri:"trainingPlanID" binding:"required"`
}
type TrainingPlanListQueryRequest struct {
	course.PaginationFilterForQuery
	EntryYears  string `json:"entry_years" form:"entry_years"`
	Departments string `json:"departments" form:"departments"`
	Degrees     string `json:"degrees" form:"degrees"`
	MajorName   string `json:"major_name" form:"major_name"`
	MajorCode   string `json:"major_code" form:"major_code"`
}
