package po

import (
	"context"
	"database/sql/driver"
	"jcourse_go/util"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SearchIndex string

// ref: https://gorm.io/zh_CN/docs/data_types.html

// warning: need manual migeration!
func (i *SearchIndex) Scan(value interface{}) error { return nil }
func (i SearchIndex) Value() (driver.Value, error)  { return nil, nil }
func (i SearchIndex) GormDataType() string          { return "tsvector" }

func (i SearchIndex) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "to_tsvector('simple', ?)",
		Vars: []interface{}{i},
	}
}

func toIndex(fields []string) SearchIndex {
	var sb strings.Builder
	for _, field := range fields {
		for _, segment := range util.Fenci(field) {
			sb.WriteString(segment)
			sb.WriteByte(' ')
		}
	}
	return SearchIndex(sb.String())
}

func (r *ReviewPO) ToSearchIndex() SearchIndex {
	return toIndex([]string{
		r.Comment,
	})
}

func (c *CoursePO) ToSearchIndex() SearchIndex {
	return toIndex([]string{
		c.Code, // 前缀模糊匹配更为适合
		c.Name,
		c.MainTeacherName,
		c.Department, // 不分词更为适合
	})
}
