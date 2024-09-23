package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"

	"github.com/SJTU-jCourse/password_hasher"

	"jcourse_go/constant"
	"jcourse_go/dal"
	"jcourse_go/model/po"
	"jcourse_go/repository"
	"jcourse_go/rpc"
)

func Login(ctx context.Context, email string, password string) (*po.UserPO, error) {
	passwordStore, err := password_hasher.MakeHashedPasswordStore(password)
	if err != nil {
		return nil, err
	}
	query := repository.NewUserQuery(dal.GetDBClient())
	userPO, err := query.GetUser(ctx, repository.WithEmail(email), repository.WithPassword(passwordStore))
	if err != nil || len(userPO) == 0 {
		return nil, err
	}
	return &userPO[0], nil
}

func Register(ctx context.Context, email string, password string, code string) (*po.UserPO, error) {
	storedCode, err := repository.GetVerifyCode(ctx, email)
	if err != nil {
		return nil, err
	}
	if storedCode != code {
		return nil, errors.New("verify code is wrong")
	}
	query := repository.NewUserQuery(dal.GetDBClient())
	userPOs, err := query.GetUser(ctx, repository.WithEmail(email))
	if err != nil {
		return nil, err
	}
	if len(userPOs) > 0 {
		return nil, errors.New("user exists for this email")
	}
	passwordStore, err := password_hasher.MakeHashedPasswordStore(password)
	if err != nil {
		return nil, err
	}
	userPO, err := query.CreateUser(ctx, email, passwordStore)
	if err != nil {
		return nil, err
	}
	_ = repository.ClearVerifyCodeHistory(ctx, email)
	return userPO, nil
}

func ResetPassword(ctx context.Context, email string, password string, code string) error {
	storedCode, err := repository.GetVerifyCode(ctx, email)
	if err != nil {
		return err
	}
	if storedCode != code {
		return errors.New("verify code is wrong")
	}
	query := repository.NewUserQuery(dal.GetDBClient())
	user, err := query.GetUser(ctx, repository.WithEmail(email))
	if err != nil {
		return err
	}
	if len(user) == 0 {
		return errors.New("user does not exist for this email")
	}
	passwordStore, err := password_hasher.MakeHashedPasswordStore(password)
	if err != nil {
		return err
	}
	err = query.ResetUserPassword(ctx, int64(user[0].ID), passwordStore)
	if err != nil {
		return err
	}
	_ = repository.ClearVerifyCodeHistory(ctx, email)
	return nil
}

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

func SendRegisterCodeEmail(ctx context.Context, email string) error {
	recentSent := repository.GetSendVerifyCodeHistory(ctx, email)
	if recentSent {
		return errors.New("recently sent code")
	}
	code, err := generateVerifyCode()
	if err != nil {
		return err
	}
	body := fmt.Sprintf(constant.EmailBodyVerifyCode, code) // nolint: gosimple
	err = repository.StoreVerifyCode(ctx, email, code)
	if err != nil {
		return err
	}
	err = rpc.SendMail(ctx, email, constant.EmailTitleVerifyCode, body)
	if err != nil {
		return err
	}
	err = repository.StoreSendVerifyCodeHistory(ctx, email)
	return err
}

func ValidateEmail(email string) bool {
	// 1. validate basic email format
	regex := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)

	if !regex.MatchString(email) { // nolint: gosimple
		return false
	}

	// 2. validate specific email model
	// TODO
	return true
}
