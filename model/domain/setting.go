package domain

import (
	"encoding/json"

	"jcourse_go/model/po"
)

type SiteSettings struct {
}

func (s *SiteSettings) UpdateFromItems(items []po.SettingItemPO) {
	itemMap := make(map[string]string)
	for _, item := range items {
		itemMap[item.Key] = item.Value
	}
	str, _ := json.Marshal(itemMap)
	_ = json.Unmarshal(str, s)
}
