package converter

import (
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/util"
)

func ConvertUserPointDetailItemFromPO(po po.UserPointDetailPO) model.UserPointDetailItem {
	return model.UserPointDetailItem{
		Time:        po.CreatedAt.Format(util.GoTimeLayout),
		Value:       po.Value,
		Description: po.Description,
	}
}
