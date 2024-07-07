package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/dto"
	"jcourse_go/model/po"
)

func ConvertUserPOToDomain(userPO po.UserPO) domain.User {
	return domain.User{
		ID:         int64(userPO.ID),
		Username:   userPO.Username,
		Email:      userPO.Email,
		Role:       userPO.UserRole,
		CreatedAt:  userPO.CreatedAt,
		LastSeenAt: userPO.LastSeenAt,
	}
}

func ConvertUserProfilePOToDomain(userProfile po.UserProfilePO) domain.UserProfile {
	return domain.UserProfile{
		UserID:     userProfile.UserID,
		Avatar:     userProfile.Avatar,
		Department: userProfile.Department,
		Type:       userProfile.Type,
		Major:      userProfile.Major,
		Degree:     userProfile.Degree,
		Grade:      userProfile.Grade,
	}
}

func PackUserWithProfile(user *domain.User, profilePO po.UserProfilePO) {
	if user == nil {
		return
	}
	profile := ConvertUserProfilePOToDomain(profilePO)
	user.Profile = profile
}

func ConvertUserDomainToReviewDTO(user domain.User) dto.UserInReviewDTO {
	return dto.UserInReviewDTO{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Profile.Avatar,
	}
}
