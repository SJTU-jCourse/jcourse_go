package command

import "jcourse_go/internal/domain/user"

type UserProfileService interface {
}

func NewUserProfileService(
	userProfileRepo user.UserProfileRepository,
) UserProfileService {
	return nil
}
