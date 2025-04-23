package model

import (
	"fmt"
	"strconv"
	"strings"

	"jcourse_go/pkg/password_hash"
)

const storeFormat = "%s$%d$%s$%s"

type Password string

func (p Password) ValidatePassword(password string) bool {
	parts := strings.Split(string(p), "$")
	if len(parts) != 4 {
		return false
	}

	algo := parts[0]
	iteration, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	salt := parts[2]
	expected := parts[3]

	h := password_hash.NewHasher(password_hash.Algorithm(algo), iteration, salt)

	actualHash := h.HashPassword(password)

	return actualHash == expected
}

func NewPassword(password string, hasher password_hash.Hasher) Password {
	hash := hasher.HashPassword(password)
	store := fmt.Sprintf(storeFormat, hasher.GetAlgorithm(), hasher.GetIteration(), hasher.GetSalt(), hash)
	return Password(store)
}
