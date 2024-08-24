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
		Vars: []interface{}{string(i)},
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

func (c *CoursePO) BeforeCreate(*gorm.DB) error {
	c.SearchIndex = toIndex([]string{
		c.Name,
		c.Code, // 前缀模糊匹配更为适合
		c.MainTeacherName,
		c.Department, // 不分词更为适合
	})
	return nil
}
func (c *CoursePO) BeforeSave(tx *gorm.DB) error {
	return c.BeforeCreate(tx)
}

func (t *TeacherPO) BeforeCreate(*gorm.DB) error {
	t.SearchIndex = toIndex([]string{
		t.Name,
		t.Department,
		t.Code,
	})
	return nil
}
func (t *TeacherPO) BeforeSave(tx *gorm.DB) error {
	return t.BeforeCreate(tx)
}

func (t *TrainingPlanPO) BeforeCreate(*gorm.DB) error {
	t.SearchIndex = toIndex([]string{
		t.Major,
		t.Department,
	})
	return nil
}
func (t *TrainingPlanPO) BeforeSave(tx *gorm.DB) error {
	return t.BeforeCreate(tx)
}

func (r *ReviewPO) BeforeCreate(*gorm.DB) error {
	r.SearchIndex = toIndex([]string{
		r.Comment,
	})
	return nil
}
func (r *ReviewPO) BeforeSave(tx *gorm.DB) error {
	return r.BeforeCreate(tx)
}
