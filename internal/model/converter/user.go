package converter

import (
	"jcourse_go/internal/domain/user"
	"jcourse_go/internal/infrastructure/entity"
	"jcourse_go/internal/interface/dto"
)

func ConvertUserDetailFromPO(po *entity.UserPO) user.UserDetail {
	return user.UserDetail{
		UserMinimal: ConvertUserMinimalFromPO(po),
		Bio:         po.Bio,
		Email:       po.Email,
		Type:        po.Type,
		Role:        po.UserRole,
		// Department:  po.Department,
		// Major:       po.Major,
		// Score:       po.Score,
		Points: po.Points,
	}
}

func ConvertUserMinimalFromPO(po *entity.UserPO) user.UserMinimal {
	return user.UserMinimal{
		ID:       po.ID,
		Username: po.Username,
		Avatar:   po.Avatar,
	}
}

func ConvertUserProfileToPO(dto dto.UserProfileDTO) entity.UserPO {
	return entity.UserPO{
		Username: dto.Username,
		Type:     dto.Type,
		Avatar:   dto.Avatar,
		// Department: dto.Department,
		// Major:      dto.Major,
		// Degree:     dto.Degree,
		// Score:      dto.Score,
		Bio: dto.Bio,
	}
}

func RemoveUserEmail(u *user.UserDetail, userID int64) {
	if u.ID == userID {
		return
	}
	u.Email = ""
}
