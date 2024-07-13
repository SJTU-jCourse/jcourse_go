package dto

type TeacherDTO struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Title      string `json:"title"`
}
