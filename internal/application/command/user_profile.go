package command

import (
	"context"

	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/domain/user"
)

type UserProfileService interface {
	UpdateUserInfo(ctx context.Context, reqCtx shared.RequestCtx, cmd user.UpdateUserInfoCommand) error
}

type userProfileService struct {
	userProfileRepo user.UserProfileRepository
}

func (u *userProfileService) UpdateUserInfo(ctx context.Context, reqCtx shared.RequestCtx, cmd user.UpdateUserInfoCommand) error {
	return u.userProfileRepo.UpdateUserInfo(ctx, reqCtx.User.UserID, cmd)
}

func NewUserProfileService(
	userProfileRepo user.UserProfileRepository,
) UserProfileService {
	return &userProfileService{
		userProfileRepo: userProfileRepo,
	}
}
