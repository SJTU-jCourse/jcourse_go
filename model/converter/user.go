package converter

import (
	"jcourse_go/model/dto"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func ConvertUserDetailFromPO(po po.UserPO) model.UserDetail {
	return model.UserDetail{
		UserMinimal: ConvertUserMinimalFromPO(po),
		Bio:         po.Bio,
		Email:       po.Email,
		Type:        po.Type,
		Role:        po.UserRole,
		Department:  po.Department,
		Major:       po.Major,
		Grade:       po.Grade,
	}
}

func ConvertUserMinimalFromPO(po po.UserPO) model.UserMinimal {
	return model.UserMinimal{
		ID:       int64(po.ID),
		Username: po.Username,
		Avatar:   po.Avatar,
	}
}

func ConvertUserProfileToPO(dto dto.UserProfileDTO) po.UserPO {
	return po.UserPO{
		Username:   dto.Username,
		Type:       dto.Type,
		Avatar:     dto.Avatar,
		Department: dto.Department,
		Major:      dto.Major,
		Degree:     dto.Degree,
		Grade:      dto.Grade,
		Bio:        dto.Bio,
	}
}

func RemoveUserEmail(u *model.UserDetail, userID int64) {
	if u.ID == userID {
		return
	}
	u.Email = ""
}
