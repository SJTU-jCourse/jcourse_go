package converter

import (
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/util"
)

func ConvertUserPointDetailItemFromPO(po po.UserPointDetailPO) model.UserPointDetailItem {
	location := util.GetLocation()
	return model.UserPointDetailItem{
		Time:        po.CreatedAt.In(location).Format(util.GoTimeLayout),
		Value:       po.Value,
		Description: po.Description,
	}
}
