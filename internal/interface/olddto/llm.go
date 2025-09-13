package olddto

type OptCourseReviewRequest struct {
	CourseName    string `json:"course_name" binding:"required"`
	ReviewContent string `json:"review_content" binding:"required"`
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
