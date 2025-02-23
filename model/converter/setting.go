package converter

import (
	"errors"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/model/types"
)

func GetSettingFromPO(po po.SettingPO) (s model.Setting, err error) {
	switch po.Type {
	case string(types.SettingTypeString):
		s = &model.StringSetting{}
	case string(types.SettingTypeInt):
		s = &model.IntSetting{}
	case string(types.SettingTypeBool):
		s = &model.BoolSetting{}
	}
	if s == nil {
		return nil, errors.New("unknown setting types")
	}
	err = s.FromPO(po)
	return
}
