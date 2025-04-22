package password

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

type PBKDF2SHA256PasswordHasher struct {
	salt      string
	iteration int
}

func (h *PBKDF2SHA256PasswordHasher) GetAlgorithm() Algorithm {
	return AlgorithmPBKDF2SHA256
}

func (h *PBKDF2SHA256PasswordHasher) HashPassword(password string) (string, error) {
	hash, _ := h.rawHash(password, h.iteration, h.salt)
	store := fmt.Sprintf(storeFormat, h.GetAlgorithm(), h.iteration, h.salt, hash)
	return store, nil
}

func (h *PBKDF2SHA256PasswordHasher) rawHash(password string, iteration int, salt string) (string, error) {
	hashed := pbkdf2.Key([]byte(password), []byte(salt), iteration, sha256.Size, sha256.New)
	return base64.StdEncoding.EncodeToString(hashed), nil
}

func (h *PBKDF2SHA256PasswordHasher) ValidatePassword(password string, hash string) bool {
	var algo, salt, expected string
	var iteration int
	_, err := fmt.Sscanf(hash, storeFormat, &algo, &iteration, &salt, &expected)
	if err != nil {
		return false
	}
	actualHash, err := h.rawHash(password, iteration, salt)
	if err != nil {
		return false
	}
	return actualHash == expected
}

func NewPBK2DFSHA256Hasher(salt string, iteration int) *PBKDF2SHA256PasswordHasher {
	return &PBKDF2SHA256PasswordHasher{
		salt:      salt,
		iteration: iteration,
	}
}
