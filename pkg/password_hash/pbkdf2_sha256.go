package password_hash

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

type PBKDF2SHA256PasswordHasher struct {
	salt      string
	iteration int
}

func (h *PBKDF2SHA256PasswordHasher) GetIteration() int {
	return h.iteration
}

func (h *PBKDF2SHA256PasswordHasher) GetSalt() string {
	return h.salt
}

func (h *PBKDF2SHA256PasswordHasher) GetAlgorithm() Algorithm {
	return AlgorithmPBKDF2SHA256
}

func (h *PBKDF2SHA256PasswordHasher) HashPassword(password string) string {
	hashed := pbkdf2.Key([]byte(password), []byte(h.salt), h.iteration, sha256.Size, sha256.New)
	return base64.StdEncoding.EncodeToString(hashed)
}

func NewPBK2DFSHA256Hasher(salt string, iteration int) Hasher {
	return &PBKDF2SHA256PasswordHasher{
		salt:      salt,
		iteration: iteration,
	}
}
