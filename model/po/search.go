package po

import (
	"context"
	"database/sql/driver"
	"strings"

	"jcourse_go/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SearchIndex string

// ref: https://gorm.io/zh_CN/docs/data_types.html

func (i *SearchIndex) Scan(value interface{}) error { return nil }
func (i *SearchIndex) Value() (driver.Value, error) { return nil, nil }
func (i *SearchIndex) GormDataType() string         { return "tsvector" }

func (i *SearchIndex) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "to_tsvector('simple', ?)",
		Vars: []interface{}{string(*i)},
	}
}

func toIndex(fields []string) SearchIndex {
	var sb strings.Builder
	for _, field := range fields {
		for _, segment := range util.SegWord(field) {
			sb.WriteString(segment)
			sb.WriteByte(' ')
		}
	}
	return SearchIndex(sb.String())
}

// ref: https://gorm.io/zh_CN/docs/hooks.html

func (c *CoursePO) BeforeCreate(*gorm.DB) error {
	if c.SearchIndex == "" {
		c.SearchIndex = toIndex([]string{
			c.Name,
			c.Code, // 前缀模糊匹配更为适合
			c.MainTeacherName,
			c.Department, // 不分词更为适合
		})
	}
	return nil
}
func (c *CoursePO) BeforeSave(tx *gorm.DB) error {
	return c.BeforeCreate(tx)
}

func (t *TeacherPO) BeforeCreate(*gorm.DB) error {
	if t.SearchIndex == "" {
		t.SearchIndex = toIndex([]string{
			t.Name,
			t.Department,
			t.Code,
		})
	}
	return nil
}
func (t *TeacherPO) BeforeSave(tx *gorm.DB) error {
	return t.BeforeCreate(tx)
}

func (t *TrainingPlanPO) BeforeCreate(*gorm.DB) error {
	if t.SearchIndex == "" {
		t.SearchIndex = toIndex([]string{
			t.Major,
			t.Department,
		})
	}
	return nil
}
func (t *TrainingPlanPO) BeforeSave(tx *gorm.DB) error {
	return t.BeforeCreate(tx)
}

func (r *ReviewPO) BeforeCreate(*gorm.DB) error {
	if r.SearchIndex == "" {
		r.SearchIndex = toIndex([]string{
			r.Comment,
		})
	}
	return nil
}
func (r *ReviewPO) BeforeSave(tx *gorm.DB) error {
	return r.BeforeCreate(tx)
}
