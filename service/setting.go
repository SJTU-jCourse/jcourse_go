package service

import (
	"context"

	"jcourse_go/internal/infra/query"
	"jcourse_go/model/converter"
	"jcourse_go/model/model"
)

func GetClientSettings(ctx context.Context) (map[string]any, error) {
	res := make(map[string]any)

	q := query.Q.SettingPO
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

func SetSettingValue(ctx context.Context, userID int64, setting model.Setting) error {
	q := query.Q.SettingPO

	po := setting.ToPO()
	po.UpdatedBy = userID

	return q.WithContext(ctx).Create(&po)
}

func GetSetting(ctx context.Context, key string) (model.Setting, error) {
	q := query.Q.SettingPO
	setting, err := q.WithContext(ctx).Where(q.Key.Eq(key)).Take()
	if err != nil {
		return nil, err
	}
	return converter.GetSettingFromPO(*setting)
}
