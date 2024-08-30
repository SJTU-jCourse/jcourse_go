package dto

type VectorizeCourseReviewsRequest struct {
	CourseID int64 `uri:"courseID" binding:"required"`
}

type GetMatchCourseRequest struct {
	Description string `json:"description" binding:"required"`
}
