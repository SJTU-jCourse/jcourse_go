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
		UserID:            userProfile.UserID,
		Avatar:            userProfile.Avatar,
		Department:        userProfile.Department,
		Type:              userProfile.Type,
		Major:             userProfile.Major,
		Degree:            userProfile.Degree,
		Grade:             userProfile.Grade,
		PersonalSignature: userProfile.PersonalSignature,
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

func ConvertToUserSummaryDTO(userPO *po.UserPO, userProfilePO *po.UserProfilePO) *dto.UserSummaryDTO {
	if userPO == nil {
		return nil
	}
	return &dto.UserSummaryDTO{
		ID:       int64(userPO.ID),
		Username: userPO.Username,
		Avatar:   userProfilePO.Avatar,
		Role:     userPO.UserRole,
	}
}

func ConvertToUserDetailsDTO(userPO *po.UserPO, userProfilePO *po.UserProfilePO) *dto.UserDetailsDTO {
	if userPO == nil {
		return nil
	}
	return &dto.UserDetailsDTO{
		ID:                int64(userPO.ID),
		Username:          userPO.Username,
		Role:              userPO.UserRole,
		LastSeenAt:        userPO.LastSeenAt,
		Type:              userPO.UserRole,
		Avatar:            userProfilePO.Avatar,
		PersonalSignature: userProfilePO.PersonalSignature,
	}
}

func ConvertToUserProfileDTO(userPO *po.UserPO, userProfilePO *po.UserProfilePO) *dto.UserProfileDTO {
	if userPO == nil {
		return nil
	}
	return &dto.UserProfileDTO{
		ID:                int64(userPO.ID),
		UserID:            userProfilePO.UserID,
		Avatar:            userProfilePO.Avatar,
		Department:        userProfilePO.Department,
		Type:              userProfilePO.Type,
		Major:             userProfilePO.Major,
		Degree:            userProfilePO.Degree,
		Grade:             userProfilePO.Grade,
		PersonalSignature: userProfilePO.PersonalSignature,
		Username:          userPO.Username,
		Role:              userPO.UserRole,
	}
}
