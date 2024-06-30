package dto

type SendEmailCodeRequest struct {
	Email string `json:"email" binding:"required"`
}

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordRequest = RegisterUserRequest

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
