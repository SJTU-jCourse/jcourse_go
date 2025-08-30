package converter

import (
	"jcourse_go/internal/model/dto"
	"jcourse_go/internal/model/model"
	"jcourse_go/internal/model/po"
)

func ConvertUserDetailFromPO(po *po.UserPO) model.UserDetail {
	return model.UserDetail{
		UserMinimal: ConvertUserMinimalFromPO(po),
		Bio:         po.Bio,
		Email:       po.Email,
		Type:        po.Type,
		Role:        po.UserRole,
		// Department:  po.Department,
		// Major:       po.Major,
		// Grade:       po.Grade,
		Points: po.Points,
	}
}

func ConvertUserMinimalFromPO(po *po.UserPO) model.UserMinimal {
	return model.UserMinimal{
		ID:       po.ID,
		Username: po.Username,
		Avatar:   po.Avatar,
	}
}

func ConvertUserProfileToPO(dto dto.UserProfileDTO) po.UserPO {
	return po.UserPO{
		Username: dto.Username,
		Type:     dto.Type,
		Avatar:   dto.Avatar,
		// Department: dto.Department,
		// Major:      dto.Major,
		// Degree:     dto.Degree,
		// Grade:      dto.Grade,
		Bio: dto.Bio,
	}
}

func RemoveUserEmail(u *model.UserDetail, userID int64) {
	if u.ID == userID {
		return
	}
	u.Email = ""
}
