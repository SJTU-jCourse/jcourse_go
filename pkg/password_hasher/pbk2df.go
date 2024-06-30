package password_hasher

import (
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/pbkdf2"
)

type PBK2DFSHA256PasswordHasher struct {
	salt      string
	iteration int64
}

func (h *PBK2DFSHA256PasswordHasher) HashPassword(password string, salt string, iteration int64) (string, error) {
	hashed := pbkdf2.Key([]byte(password), []byte(salt), int(iteration), sha256.Size, sha256.New)
	res := base64.StdEncoding.EncodeToString(hashed)
	return res, nil
}
