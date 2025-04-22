package validator

import (
	"fmt"
	"testing"
)

func TestNewEmailValidator(t *testing.T) {
	tests := []struct {
		name         string
		suffixDomain string
		expectedType string
	}{
		{"Empty suffixDomain", "", "*validator.CommonEmailValidator"},
		{"Non-empty suffixDomain", "example.com", "*validator.SuffixEmailValidator"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewEmailValidator(tt.suffixDomain)
			if gotType := fmt.Sprintf("%T", validator); gotType != tt.expectedType {
				t.Errorf("expected type %s, got %s", tt.expectedType, gotType)
			}
		})
	}
}

func TestCommonEmailValidator_Validate(t *testing.T) {
	validator := &CommonEmailValidator{}
	tests := []struct {
		name   string
		email  string
		result bool
	}{
		{"Valid email", "user@example.com", true},
		{"Valid email with subdomain", "user@mail.example.com", true},
		{"Invalid email - no domain", "user@", false},
		{"Invalid email - no username", "@example.com", false},
		{"Invalid email - no @", "userexample.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validator.Validate(tt.email); got != tt.result {
				t.Errorf("expected %v, got %v", tt.result, got)
			}
		})
	}
}

func TestSuffixEmailValidator_Validate(t *testing.T) {
	validator := &SuffixEmailValidator{suffixDomain: "example.com"}
	tests := []struct {
		name   string
		email  string
		result bool
	}{
		{"Valid email with suffix", "user@example.com", true},
		{"Valid email with subdomain in suffix", "user@mail.example.com", false},
		{"Invalid email - no suffix", "user@other.com", false},
		{"Invalid email - no username", "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validator.Validate(tt.email); got != tt.result {
				t.Errorf("expected %v, got %v", tt.result, got)
			}
		})
	}
}
