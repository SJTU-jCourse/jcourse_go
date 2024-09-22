package service

import (
	"context"

	"jcourse_go/dal"
	"jcourse_go/repository"
)

func GetClientSettings(ctx context.Context) (map[string]any, error) {
	res := make(map[string]any)
	query := repository.NewSettingQuery(dal.GetDBClient())
	settings, err := query.GetClientSettings(ctx)
	if err != nil {
		return nil, err
	}
	for _, setting := range settings {
		res[setting.GetKey()] = setting.GetValue()
	}
	return res, nil
}
