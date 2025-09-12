package converter

import (
	"jcourse_go/internal/domain/user"
	"jcourse_go/internal/infrastructure/entity"
	"jcourse_go/pkg/util"
)

func ConvertUserPointDetailItemFromPO(po entity.UserPointDetail) user.UserPointDetailItem {
	location := util.GetLocation()
	return user.UserPointDetailItem{
		Time:        po.CreatedAt.In(location).Format(util.GoTimeLayout),
		Value:       po.Value,
		Description: po.Description,
	}
}
