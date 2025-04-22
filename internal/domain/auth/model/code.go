package model

import "time"

type VerificationCode struct {
	Code      string
	Email     string
	CreatedAt time.Time
	TTL       time.Duration
}

func (v *VerificationCode) InRateLimit(limit time.Duration) bool {
	return v.CreatedAt.Add(limit).After(time.Now())
}

func (v *VerificationCode) IsExpired() bool {
	return time.Now().After(v.CreatedAt.Add(v.TTL))
}

func (v *VerificationCode) IsValid(code string) bool {
	return v.Code == code && !v.IsExpired()
}
