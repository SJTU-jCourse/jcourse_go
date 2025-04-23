package model

import (
	"testing"

	"jcourse_go/pkg/password_hash"

	"github.com/stretchr/testify/assert"
)

type mockHasher struct {
	algorithm string
	iteration int
	salt      string
	hash      string
}

func (m *mockHasher) HashPassword(password string) string {
	return m.hash
}

func (m *mockHasher) GetAlgorithm() password_hash.Algorithm {
	return password_hash.Algorithm(m.algorithm)
}

func (m *mockHasher) GetIteration() int {
	return m.iteration
}

func (m *mockHasher) GetSalt() string {
	return m.salt
}

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name          string
		inputPassword string
		hasher        *mockHasher
		expectedStore string
	}{
		{
			name:          "create new password",
			inputPassword: "password123",
			hasher: &mockHasher{
				algorithm: "alg",
				iteration: 1000,
				salt:      "randomSalt",
				hash:      "hashedPassword",
			},
			expectedStore: "alg$1000$randomSalt$hashedPassword",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPassword(tt.inputPassword, tt.hasher)
			assert.Equal(t, Password(tt.expectedStore), result)
		})
	}
}

func TestPassword_IsValid(t *testing.T) {

	tests := []struct {
		Name     string
		Password string
		Store    Password
		Want     bool
	}{
		{
			Name:     "valid password",
			Password: "test",
			Store:    Password("pbkdf2_sha256$720000$$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM="),
			Want:     true,
		},
		{Name: "not ok 1", Password: "test2", Store: "pbkdf2_sha256$720000$$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: false},
		{Name: "not ok 2", Password: "test", Store: "pbkdf2_sha256$72000$$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: false},
		{Name: "not ok 3", Password: "test", Store: "pbkdf2_sha256$720000$123$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: false},
		{Name: "not ok 4", Password: "test", Store: "pbkdf2_sha256$720000$$j2DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: false},
		{Name: "not ok 5", Password: "test", Store: "", Want: false},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Want, tt.Store.ValidatePassword(tt.Password))
		})
	}
}
