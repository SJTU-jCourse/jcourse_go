package store

import (
	"encoding/hex"

	"golang.org/x/crypto/blake2b"
)

type UsernameChanger interface {
	ChangeUsername(username string) string
}

type TrivialUsernameChanger struct{}

func (c *TrivialUsernameChanger) ChangeUsername(username string) string {
	return username
}

type HashUsernameChanger struct {
	salt string
}

func (c *HashUsernameChanger) ChangeUsername(username string) string {
	hasher, _ := blake2b.New(16, nil)
	hasher.Write([]byte(username + c.salt))
	bytes := hasher.Sum(nil)
	return hex.EncodeToString(bytes)
}
