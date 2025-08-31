package auth

type SendVerificationCodeCommand struct {
	Email string `json:"email" binding:"required"`
}

type RegisterUserCommand struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordCommand RegisterUserCommand

type LoginCommand struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
