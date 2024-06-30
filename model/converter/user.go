package converter

import (
	"jcourse_go/model/domain"
	"jcourse_go/model/po"
)

func UserPOToDomain(userPO *po.UserPO) *domain.User {
	if userPO == nil {
		return nil
	}
	return &domain.User{
		ID:         userPO.ID,
		Username:   userPO.Username,
		Email:      userPO.Email,
		Role:       userPO.UserRole,
		CreatedAt:  userPO.CreatedAt,
		LastSeenAt: userPO.LastSeenAt,
	}
}
