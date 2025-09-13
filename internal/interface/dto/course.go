package dto

type CourseDetailRequest struct {
	CourseID int64 `uri:"courseID" binding:"required"`
}
