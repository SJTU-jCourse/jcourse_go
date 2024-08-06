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
		Bio:        userProfile.Bio,
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

func ConvertToUserDetailDTO(userPO *po.UserPO, userProfilePO *po.UserProfilePO) *dto.UserDetailDTO {
	if userPO == nil {
		return nil
	}
	return &dto.UserDetailDTO{
		ID:       int64(userPO.ID),
		Username: userPO.Username,
		Avatar:   userProfilePO.Avatar,
		Bio:      userProfilePO.Bio,
	}
}

func ConvertToUserProfileDTO(userPO *po.UserPO, userProfilePO *po.UserProfilePO) *dto.UserProfileDTO {
	if userPO == nil {
		return nil
	}
	return &dto.UserProfileDTO{
		ID:         int64(userPO.ID),
		UserID:     userProfilePO.UserID,
		Username:   userPO.Username,
		Bio:        userProfilePO.Bio,
		Email:      userPO.Email,
		Avatar:     userProfilePO.Avatar,
		Role:       userPO.UserRole,
		Department: userProfilePO.Department,
		Major:      userProfilePO.Major,
		Grade:      userProfilePO.Grade,
	}
}
