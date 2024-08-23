package repository

import (
	"jcourse_go/util"
	"strings"

	"gorm.io/gorm"
)

func (c *ReviewQuery) WithSearch(query string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("search_idx @@ to_tsquery('simple', ?)",
			strings.Join(util.Fenci(query), " | "),
		)
	}
}
