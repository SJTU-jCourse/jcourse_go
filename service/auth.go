package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"

	"jcourse_go/constant"
	"jcourse_go/model/domain"
	"jcourse_go/pkg/password_hasher"
	"jcourse_go/repository"
	"jcourse_go/rpc"
)

func Login(ctx context.Context, email string, password string) (*domain.User, error) {
	passwordStore, err := password_hasher.MakeHashedPasswordStore(password)
	if err != nil {
		return nil, err
	}
	query := repository.NewUserQuery()
	user, err := query.GetUserDetail(ctx, query.WithEmail(email), query.WithPassword(passwordStore))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func Register(ctx context.Context, email string, password string, code string) (*domain.User, error) {
	storedCode, err := repository.GetVerifyCode(ctx, email)
	if err != nil {
		return nil, err
	}
	if storedCode != code {
		return nil, errors.New("verify code is wrong")
	}
	query := repository.NewUserQuery()
	user, err := query.GetUserDetail(ctx, query.WithEmail(email))
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("user exists for this email")
	}
	passwordStore, err := password_hasher.MakeHashedPasswordStore(password)
	if err != nil {
		return nil, err
	}
	user, err = query.CreateUser(ctx, email, passwordStore)
	if err != nil {
		return nil, err
	}
	_ = repository.ClearVerifyCodeHistory(ctx, email)
	return user, nil
}

func ResetPassword(ctx context.Context, email string, password string, code string) error {
	storedCode, err := repository.GetVerifyCode(ctx, email)
	if err != nil {
		return err
	}
	if storedCode != code {
		return errors.New("verify code is wrong")
	}
	query := repository.NewUserQuery()
	user, err := query.GetUserDetail(ctx, query.WithEmail(email))
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user does not exist for this email")
	}
	passwordStore, err := password_hasher.MakeHashedPasswordStore(password)
	if err != nil {
		return err
	}
	err = query.ResetUserPassword(ctx, user, passwordStore)
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
	body := fmt.Sprintf(constant.EmailBodyVerifyCode, code)
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
	if !regex.MatchString(email) {
		return false
	}

	// 2. validate specific email domain
	// TODO
	return true
}
