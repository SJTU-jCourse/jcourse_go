package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeHashedPasswordStore(t *testing.T) {
	hasher := NewPBK2DFSHA256Hasher("", 720000)
	store, err := hasher.HashPassword("test")
	assert.Nil(t, err)
	assert.Equal(t, "pbkdf2_sha256$720000$$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", store)
}

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		Name     string
		Password string
		Store    string
		Want     bool
	}{
		{Name: "ok", Password: "test", Store: "pbkdf2_sha256$720000$$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: true},
		{Name: "not ok", Password: "test2", Store: "pbkdf2_sha256$720000$$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: false},
		{Name: "not ok", Password: "test", Store: "pbkdf2_sha256$72000$$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: false},
		{Name: "not ok", Password: "test", Store: "pbkdf2_sha256$720000$123$j92DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: false},
		{Name: "not ok", Password: "test", Store: "pbkdf2_sha256$720000$$j2DYtoNAOL6uFf22YbNOKRJo8Q9a2gjXij5KkKhLPM=", Want: false},
		{Name: "not ok", Password: "test", Store: "", Want: false},
	}

	hasher := NewPBK2DFSHA256Hasher("", 720000)
	for _, testcase := range cases {
		t.Run(testcase.Name, func(t *testing.T) {
			ok := hasher.ValidatePassword(testcase.Password, testcase.Store)
			assert.Equal(t, testcase.Want, ok)
		})
	}

}
