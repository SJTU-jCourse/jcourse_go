package dto

type OptCourseReviewRequest struct {
	CourseName    string `json:"courseName" binding:"required"`
	ReviewContent string `json:"reviewContent" binding:"required"`
}

type OptCourseReviewResponse struct {
	Suggestion *string `json:"suggestion"`
	Result     *string `json:"result"`
}

type GetCourseSummaryResponse struct {
	Summary *string `json:"summary"`
}

type GetMatchCourseRequest struct {
	Description string `json:"description" binding:"required"`
}
