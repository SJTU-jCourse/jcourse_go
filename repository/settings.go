package repository

import (
	"context"

	"gorm.io/gorm"

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
	db := dal.GetDBClient().WithContext(ctx).Model(&po.SettingItemPO{})
	// 无历史记录，只写入
	var count int64
	db.Where(&po.SettingItemPO{Key: key}).Count(&count)
	if count == 0 {
		err := db.Create(&po.SettingItemPO{Key: key, Value: value, UpdatedBy: userID}).Error
		if err != nil {
			return err
		}
		return nil
	}
	// 有历史记录，需要删除旧的，写入新的（同一个事务）
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("key = ?", key).Delete(&po.SettingItemPO{}).Error
		if err != nil {
			return err
		}
		err = tx.Create(&po.SettingItemPO{Key: key, Value: value, UpdatedBy: userID}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
