package converter

import (
	"time"

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

func ConvertUserDomainToUserSummaryDTO(id int64, reviewCount int64, likeReceive int64, tipReceive int64, followingCourseCount int64) *dto.UserSummaryDTO {
	return &dto.UserSummaryDTO{
		ID:                   id,
		ReviewCount:          reviewCount,
		LikeReceive:          likeReceive,
		TipReceive:           tipReceive,
		FollowingCourseCount: followingCourseCount,
	}
}

func ConvertUserDomainToUserDetailDTO(userDomain *domain.User) *dto.UserDetailDTO {
	if userDomain == nil {
		return nil
	}
	return &dto.UserDetailDTO{
		ID:       userDomain.ID,
		Username: userDomain.Username,
		Avatar:   userDomain.Profile.Avatar,
		Bio:      userDomain.Profile.Bio,
	}
}

func ConvertUserDomainToUserProfileDTO(userDomain *domain.User) *dto.UserProfileDTO {
	if userDomain == nil {
		return nil
	}
	return &dto.UserProfileDTO{
		UserID:     userDomain.ID,
		Username:   userDomain.Username,
		Bio:        userDomain.Profile.Bio,
		Email:      userDomain.Email,
		Avatar:     userDomain.Profile.Avatar,
		Role:       userDomain.Role,
		Department: userDomain.Profile.Department,
		Major:      userDomain.Profile.Major,
		Grade:      userDomain.Profile.Grade,
	}
}

func ConvertUpdateUserProfileDTOToUserPO(userProfileDTO *dto.UserProfileDTO, userPO *po.UserPO) po.UserPO {
	updatedUserPO := po.UserPO{
		Username:   userProfileDTO.Username,
		Email:      userPO.Email,
		Password:   userPO.Password,
		UserRole:   userPO.UserRole,
		LastSeenAt: time.Now(),
	}
	if userProfileDTO.UserID != 0 {
		updatedUserPO.ID = uint(userProfileDTO.UserID)
	}
	return updatedUserPO
}

func ConvertUpdateUserProfileDTOToUsrProfilePO(userProfileDTO *dto.UserProfileDTO, userProfilePO *po.UserProfilePO) po.UserProfilePO {
	// 保留一些immutable的属性
	updatedUserProfilePO := po.UserProfilePO{
		UserID:     userProfilePO.UserID,
		Avatar:     userProfileDTO.Avatar,
		Department: userProfileDTO.Department,
		Type:       userProfilePO.Type,
		Major:      userProfileDTO.Major,
		Degree:     userProfilePO.Degree,
		Grade:      userProfileDTO.Grade,
		Bio:        userProfileDTO.Bio,
	}
	if userProfileDTO.UserID != 0 {
		updatedUserProfilePO.ID = uint(userProfileDTO.UserID)
	}
	return updatedUserProfilePO
}
