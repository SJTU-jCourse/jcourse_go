package dto

type TeacherDTO struct {
	ID         int64           `json:"id"`
	Email      string          `json:"email"`
	Code       string          `json:"code"`
	Name       string          `json:"name"`
	Department string          `json:"department"`
	Title      string          `json:"title"`
	Picture    string          `json:"picture"`
	ProfileURL string          `json:"profileURL"`
	Courses    []BaseCourseDTO `json:"courses"`
}

type TeacherDetailRequest struct {
	TeacherID int64 `uri:"teacherID" binding:"required"`
}

type TeacherDetailResponse = TeacherDTO

type TeacherQueryRequest struct {
	Text string `json:"text"`
	Department string `json:"department"`
	MajorName string `json:"major_name"`
	MajorCode string `json:"major_code"`
	SortDirection string `json:"sort_direction"`
	SortBy string `json:"sort_by"`
}
