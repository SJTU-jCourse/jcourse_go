package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/SJTU-jCourse/password_hasher"

	"jcourse_go/constant"
	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/repository"
	"jcourse_go/rpc"
)

func Login(ctx context.Context, email string, password string) (*model.UserDetail, error) {
	emailToQuery := convertEmailToQuery(email)

	u := repository.Q.UserPO
	userPO, err := u.WithContext(ctx).Where(u.Email.Eq(emailToQuery)).Limit(1).Take()
	if err != nil {
		return nil, err
	}
	if userPO == nil {
		return nil, errors.New("user does not exist for this email")
	}
	ok, err := password_hasher.ValidatePassword(password, userPO.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("password is wrong")
	}
	user := converter.ConvertUserDetailFromPO(userPO)
	return &user, nil
}

func Register(ctx context.Context, email string, password string, code string) (*model.UserDetail, error) {
	storedCode, err := repository.GetVerifyCode(ctx, email)
	if err != nil {
		return nil, err
	}
	if storedCode != code {
		return nil, errors.New("verify code is wrong")
	}

	emailToQuery := convertEmailToQuery(email)

	u := repository.Q.UserPO
	userPO, err := u.WithContext(ctx).Where(u.Email.Eq(emailToQuery)).Limit(1).Take()
	if err != nil {
		return nil, err
	}
	if userPO != nil {
		return nil, errors.New("user exists for this email")
	}
	passwordStore, err := password_hasher.MakeHashedPasswordStore(password)
	if err != nil {
		return nil, err
	}

	user := buildUserToCreate(emailToQuery, passwordStore)
	err = u.WithContext(ctx).Create(&user)
	if err != nil {
		return nil, err
	}

	_ = repository.ClearVerifyCodeHistory(ctx, email)

	userDetail := converter.ConvertUserDetailFromPO(&user)
	return &userDetail, nil
}

func ResetPassword(ctx context.Context, email string, password string, code string) error {
	storedCode, err := repository.GetVerifyCode(ctx, email)
	if err != nil {
		return err
	}
	if storedCode != code {
		return errors.New("verify code is wrong")
	}

	emailToQuery := convertEmailToQuery(email)

	u := repository.Q.UserPO
	userPO, err := u.WithContext(ctx).Where(u.Email.Eq(emailToQuery)).Limit(1).Take()
	if err != nil {
		return err
	}
	if userPO == nil {
		return errors.New("user does not exist for this email")
	}

	passwordStore, err := password_hasher.MakeHashedPasswordStore(password)
	if err != nil {
		return err
	}

	_, err = u.WithContext(ctx).Select(u.Password).Where(u.ID.Eq(userPO.ID)).Update(u.Password, passwordStore)
	if err != nil {
		return err
	}
	_ = repository.ClearVerifyCodeHistory(ctx, email)
	return nil
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
		fmt.Printf("StoreVerifyCode error: %v\n", err)
		return err
	}

	sender := rpc.SMTPSender{}
	err = sender.SendMail(ctx, email, constant.EmailTitleVerifyCode, body)
	if err != nil {
		fmt.Printf("SendMail error: %v\n", err)
		return err
	}
	err = repository.StoreSendVerifyCodeHistory(ctx, email)
	return err
}

func SendRegisterCodeEmailMock(ctx context.Context, email string) error {
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
		fmt.Printf("StoreVerifyCode error: %v\n", err)
		return err
	}
	fmt.Printf("[HINT] Send email to %s, title: %s, body: %s\n", email, constant.EmailTitleVerifyCode, body)
	err = repository.StoreSendVerifyCodeHistory(ctx, email)
	return err
}
