package repository

import (
	"jcourse_go/util"
	"strings"

	"gorm.io/gorm"
)

func (*CourseQuery) WithSearch(query string) DBOption       { return withSearch(query) }
func (*ReviewQuery) WithSearch(query string) DBOption       { return withSearch(query) }
func (*TeacherQuery) WithSearch(query string) DBOption      { return withSearch(query) }
func (*TrainingPlanQuery) WithSearch(query string) DBOption { return withSearch(query) }

func withSearch(query string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("search_index @@ to_tsquery('simple', ?)",
			userQueryToTsQuery(query),
		)
	}
}

// 目前只搜用户名
func (*UserQuery) WithSearch(query string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name LIKE ?", query+"%")
	}
}

func userQueryToTsQuery(query string) string {

	return strings.Join(util.Fenci(query), " | ")
}
