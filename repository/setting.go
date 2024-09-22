package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/model/model"
	"jcourse_go/model/po"
)

func GetSettingFromPO(po po.SettingPO) (s model.Setting, err error) {
	switch po.Type {
	case model.SettingTypeString:
		s = &model.StringSetting{}
	case model.SettingTypeInt:
		s = &model.IntSetting{}
	case model.SettingTypeBool:
		s = &model.BoolSetting{}
	}
	if s == nil {
		return nil, errors.New("unknown setting type")
	}
	err = s.FromPO(po)
	return
}

type ISettingQuery interface {
	GetSetting(ctx context.Context, key string) (model.Setting, error)
	SetSetting(ctx context.Context, userID int64, setting model.Setting) error
	GetClientSettings(ctx context.Context) ([]model.Setting, error)
}

type SettingQuery struct {
	db *gorm.DB
}

func (s *SettingQuery) GetClientSettings(ctx context.Context) ([]model.Setting, error) {
	var settingPOs []po.SettingPO
	err := s.db.WithContext(ctx).Model(&po.SettingPO{}).Where("client = true").Find(&settingPOs).Error
	if err != nil {
		return nil, err
	}
	res := make([]model.Setting, 0, len(settingPOs))
	for _, settingPO := range settingPOs {
		setting, err := GetSettingFromPO(settingPO)
		if err != nil {
			continue
		}
		res = append(res, setting)
	}
	return res, nil
}

func (s *SettingQuery) SetSetting(ctx context.Context, userID int64, setting model.Setting) error {
	settingPO := setting.ToPO()
	settingPO.UpdatedBy = userID
	err := s.db.WithContext(ctx).Model(&po.SettingPO{}).
		Clauses(clause.OnConflict{UpdateAll: true}).Create(&settingPO).Error
	return err
}

func (s *SettingQuery) GetSetting(ctx context.Context, key string) (model.Setting, error) {
	res := po.SettingPO{}
	err := s.db.WithContext(ctx).Model(&po.SettingPO{}).Where("key = ?", key).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return GetSettingFromPO(res)
}

func NewSettingQuery(db *gorm.DB) ISettingQuery {
	return &SettingQuery{db: db}
}
