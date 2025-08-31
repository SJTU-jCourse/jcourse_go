package auth

import "time"

type VerificationCode struct {
	Email        string
	Code         string
	ValidateTill time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
