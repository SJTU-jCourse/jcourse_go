package model

import (
	"time"

	"jcourse_go/model/types"
	"jcourse_go/pkg/password"
	"jcourse_go/pkg/store"
)

type UserDomain struct {
	ID       int
	Username string // hash in need
	Email    string // store in need
	Password string // hashed

	Role          types.UserRole
	SuspendedTill *time.Time
}

func (u *UserDomain) ValidatePassword(password string, validator password.Validator) bool {
	return validator.ValidatePassword(password, u.Password)
}

func (u *UserDomain) SetPassword(password string, hasher password.Hasher) (err error) {
	u.Password, err = hasher.HashPassword(password)
	return
}

func (u *UserDomain) SetUsername(username string, changer store.UsernameChanger) (err error) {
	u.Username = changer.ChangeUsername(username)
	return nil
}

func (u *UserDomain) CanLogin() bool {
	return u.SuspendedTill == nil || u.SuspendedTill.After(time.Now())
}

func (u *UserDomain) IsAdmin() bool {
	return u.Role == types.UserRoleAdmin
}
