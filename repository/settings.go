package repository

import (
	"context"

	"jcourse_go/dal"
	"jcourse_go/model/po"
)

const RedisKeySiteSetting = "site_settings"

func getSiteSetting(ctx context.Context, key string) (string, error) {
	client := dal.GetRedisClient()
	val, err := client.HGet(ctx, RedisKeySiteSetting, key).Result()
	// 缓存中存在
	if err == nil {
		return val, nil
	}

	// 缓存不存在，查询 db
	db := dal.GetDBClient().WithContext(ctx).Model(&po.SettingItemPO{})
	settingItem := &po.SettingItemPO{}
	err = db.Where("key = ?", key).First(settingItem).Error
	if err != nil {
		return "", err
	}
	// 写入缓存
	_, _ = client.HSet(ctx, RedisKeySiteSetting, key, settingItem.Value).Result()
	return settingItem.Value, nil
}

func GetSiteSetting(ctx context.Context, key string) (string, error) {
	return getSiteSetting(ctx, key)
}

func SetSiteSetting(ctx context.Context, key string, value string, userID int64) error {

	// 无历史记录，只写入

	// 有历史记录，需要删除旧的，写入新的（同一个事务）

	return nil
}
