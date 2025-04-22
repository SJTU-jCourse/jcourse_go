package model

import (
	"testing"
	"time"
)

func TestVerificationCode_InRateLimit(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		createdAt  time.Time
		limit      time.Duration
		wantResult bool
	}{
		{"within limit", now.Add(-1 * time.Minute), 2 * time.Minute, true},
		{"outside limit", now.Add(-3 * time.Minute), 2 * time.Minute, false},
		{"on limit boundary", now.Add(-2 * time.Minute), 2 * time.Minute, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := VerificationCode{CreatedAt: tc.createdAt}
			got := v.InRateLimit(tc.limit)
			if got != tc.wantResult {
				t.Errorf("InRateLimit() = %v, want %v", got, tc.wantResult)
			}
		})
	}
}

func TestVerificationCode_IsExpired(t *testing.T) {

	now := time.Now()

	tests := []struct {
		name       string
		createdAt  time.Time
		ttl        time.Duration
		wantResult bool
	}{
		{"not expired", now.Add(-1 * time.Minute), 2 * time.Minute, false},
		{"expired", now.Add(-3 * time.Minute), 2 * time.Minute, true},
		{"on expiry boundary", now.Add(-2 * time.Minute), 2 * time.Minute, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := VerificationCode{CreatedAt: tc.createdAt, TTL: tc.ttl}
			got := v.IsExpired()
			if got != tc.wantResult {
				t.Errorf("IsExpired() = %v, want %v", got, tc.wantResult)
			}
		})
	}
}

func TestVerificationCode_IsValid(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		code       VerificationCode
		inputCode  string
		wantResult bool
	}{
		{"valid code not expired", VerificationCode{Code: "12345", TTL: 2 * time.Minute, CreatedAt: now.Add(-1 * time.Minute)}, "12345", true},
		{"expired code", VerificationCode{Code: "12345", TTL: 2 * time.Minute, CreatedAt: now.Add(-3 * time.Minute)}, "12345", false},
		{"incorrect code", VerificationCode{Code: "12345", TTL: 2 * time.Minute, CreatedAt: now.Add(-1 * time.Minute)}, "54321", false},
		{"expired and incorrect code", VerificationCode{Code: "12345", TTL: 2 * time.Minute, CreatedAt: now.Add(-3 * time.Minute)}, "54321", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.code.IsValid(tc.inputCode)
			if got != tc.wantResult {
				t.Errorf("IsValid(%q) = %v, want %v", tc.inputCode, got, tc.wantResult)
			}
		})
	}
}
