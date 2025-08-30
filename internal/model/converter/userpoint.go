package converter

import (
	"jcourse_go/internal/model/model"
	"jcourse_go/internal/model/po"
	"jcourse_go/pkg/util"
)

func ConvertUserPointDetailItemFromPO(po po.UserPointDetailPO) model.UserPointDetailItem {
	location := util.GetLocation()
	return model.UserPointDetailItem{
		Time:        po.CreatedAt.In(location).Format(util.GoTimeLayout),
		Value:       po.Value,
		Description: po.Description,
	}
}
