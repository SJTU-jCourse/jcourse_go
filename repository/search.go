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

// 空格分割的每个词都要匹配，词内分词做模糊匹配
func userQueryToTsQuery(query string) string {
	var sb strings.Builder
	words := strings.Fields(query)
	for i, word := range words {
		sb.WriteByte('(')
		segs := util.Fenci(word)
		for j, seg := range segs {
			sb.WriteString(seg)
			if j != len(segs)-1 {
				sb.WriteString(" | ")
			}
		}
		sb.WriteByte(')')
		if i != len(words) {
			sb.WriteString(" & ")
		}
	}
	return sb.String()
}
