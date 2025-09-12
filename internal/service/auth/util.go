package auth

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"

	"jcourse_go/internal/constant"
	"jcourse_go/internal/infrastructure/entity"
	"jcourse_go/internal/model/types"
	"jcourse_go/pkg/util"
)

func generateVerifyCode() (string, error) {
	var number []byte
	for i := 0; i < constant.AuthVerifyCodeLen; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		number = append(number, constant.VerifyCodeDigits[n.Int64()])
	}

	return string(number), nil
}

func hashUsername(username string, salt string) string {
	hasher, _ := blake2b.New(16, nil)
	hasher.Write([]byte(username + salt))
	bytes := hasher.Sum(nil)
	return hex.EncodeToString(bytes)
}

func extractUsername(email string) string {
	s := strings.Split(email, "@")
	return s[0]
}

func convertEmailToQuery(email string) string {
	if !util.IsEnableSJTUFeature() {
		return email
	}
	username := extractUsername(email)
	hashedUsername := hashUsername(username, util.GetHashSalt())
	return hashedUsername
}

func buildUserToCreate(email string, passwordStore string) entity.User {
	return entity.User{
		Username:   email,
		Email:      email,
		UserRole:   string(types.UserRoleNormal),
		LastSeenAt: time.Now(),
		Password:   passwordStore,
	}
}
