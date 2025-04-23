package model

import (
	"time"

	"jcourse_go/model/types"
)

type UserDomain struct {
	ID       int
	Username string   // hash in need
	Email    string   // store in need
	Password Password // hashed

	Role types.UserRole

	CreatedAt  time.Time
	LastSeenAt time.Time

	SuspendedAt   *time.Time
	SuspendedTill *time.Time
}

func (u *UserDomain) ValidatePassword(password string) bool {
	return u.Password.ValidatePassword(password)
}

func (u *UserDomain) SetPassword(password Password) {
	u.Password = password
}

func (u *UserDomain) CanLogin() bool {
	return u.SuspendedTill == nil || u.SuspendedTill.After(time.Now())
}

func (u *UserDomain) Suspend(till time.Time) {
	now := time.Now()
	u.SuspendedAt = &now
	u.SuspendedTill = &till
}

func (u *UserDomain) IsAdmin() bool {
	return u.Role == types.UserRoleAdmin
}

func (u *UserDomain) SetAdmin() {
	u.Role = types.UserRoleAdmin
}

func (u *UserDomain) Seen() {
	u.LastSeenAt = time.Now()
}

func NewUser(email string, password Password) *UserDomain {
	now := time.Now()
	return &UserDomain{
		Email:     email,
		Password:  password,
		Role:      types.UserRoleNormal,
		CreatedAt: now,
	}
}
