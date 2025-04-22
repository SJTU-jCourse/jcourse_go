package store

import (
	"testing"
)

func TestTrivialUsernameChanger_ChangeUsername(t *testing.T) {
	changer := &TrivialUsernameChanger{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"basic", "username", "username"},
		{"specialCharacters", "@user!#123", "@user!#123"},
		{"numeric", "123456", "123456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := changer.ChangeUsername(tt.input)
			if actual != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestHashUsernameChanger_ChangeUsername(t *testing.T) {
	changer := &HashUsernameChanger{salt: "saltValue"}

	tests := []struct {
		name         string
		input        string
		expectedFunc func(string) bool
	}{
		{"empty", "", func(val string) bool { return val != "" }},
		{"basic", "username", func(val string) bool { return val != "" }},
		{"specialCharacters", "@user!#123", func(val string) bool { return val != "" }},
		{"numeric", "123456", func(val string) bool { return val != "" }},
		{"sameInput", "username", func(val string) bool { return len(val) == 32 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := changer.ChangeUsername(tt.input)
			if !tt.expectedFunc(actual) {
				t.Errorf("unexpected hash result for input %s: %s", tt.input, actual)
			}
		})
	}
}
