package converter

import (
	"errors"

	"jcourse_go/internal/domain"
	"jcourse_go/internal/infrastructure/entity"
	"jcourse_go/internal/model/types"
)

func GetSettingFromPO(po entity.SettingPO) (s domain.Setting, err error) {
	switch po.Type {
	case string(types.SettingTypeString):
		s = &domain.StringSetting{}
	case string(types.SettingTypeInt):
		s = &domain.IntSetting{}
	case string(types.SettingTypeBool):
		s = &domain.BoolSetting{}
	}
	if s == nil {
		return nil, errors.New("unknown setting types")
	}
	err = s.FromPO(po)
	return
}
