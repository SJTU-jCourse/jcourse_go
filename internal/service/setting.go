package service

import (
	"context"

	"jcourse_go/internal/domain"
	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/internal/model/converter"
)

func GetClientSettings(ctx context.Context) (map[string]any, error) {
	res := make(map[string]any)

	q := repository.Q.SettingPO
	settings, err := q.WithContext(ctx).Where(q.Client.Is(true)).Find()
	if err != nil {
		return nil, err
	}

	for _, setting := range settings {
		setModel, _ := converter.GetSettingFromPO(*setting)
		res[setModel.GetKey()] = setModel.GetValue()
	}
	return res, nil
}

func SetSettingValue(ctx context.Context, userID int64, setting domain.Setting) error {
	q := repository.Q.SettingPO

	po := setting.ToPO()
	po.UpdatedBy = userID

	return q.WithContext(ctx).Create(&po)
}

func GetSetting(ctx context.Context, key string) (domain.Setting, error) {
	q := repository.Q.SettingPO
	setting, err := q.WithContext(ctx).Where(q.Key.Eq(key)).Take()
	if err != nil {
		return nil, err
	}
	return converter.GetSettingFromPO(*setting)
}
