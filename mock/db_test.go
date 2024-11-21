package mock

import (
	"jcourse_go/dal"
	"jcourse_go/model/po"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockDB(t *testing.T) {
	t.Run("TestMockDB", func(t *testing.T) {
		param := MockDBParams{
			Seed:      0,
			BatchSize: 100,
			Users:     10,
			Courses:   10,
			Reviews:   10,
			Teachers:  10,
		}
		_, err := MockDB(param)
		assert.Nil(t, err)
		db := dal.GetDBClient()
		var userCount int64
		var courseCount int64
		var reviewCount int64
		var teacherCount int64
		db.Model(&po.UserPO{}).Count(&userCount)
		db.Model(&po.CoursePO{}).Count(&courseCount)
		db.Model(&po.ReviewPO{}).Count(&reviewCount)
		db.Model(&po.TeacherPO{}).Count(&teacherCount)
		assert.Equal(t, int64(param.Users), userCount)
		assert.Equal(t, int64(param.Courses), courseCount)
		assert.Equal(t, int64(param.Reviews), reviewCount)
		assert.Equal(t, int64(param.Teachers), teacherCount)
		courses := make([]po.CoursePO, 0)
		db.Model(&po.CoursePO{}).Find(&courses)
		for _, course := range courses[:5] {
			t.Logf("Course: %v", course)
		}
	})
	t.Run("TestLarger", func(t *testing.T) {
		param := MockDBParams{
			Seed:      0,
			BatchSize: 100,
			Users:     10000,
			Courses:   20000,
			Reviews:   30000,
			Teachers:  5000,
		}
		_, err := MockDB(param)
		assert.Nil(t, err)
		db := dal.GetDBClient()
		var userCount int64
		var courseCount int64
		var reviewCount int64
		var teacherCount int64
		db.Model(&po.UserPO{}).Count(&userCount)
		db.Model(&po.CoursePO{}).Count(&courseCount)
		db.Model(&po.ReviewPO{}).Count(&reviewCount)
		db.Model(&po.TeacherPO{}).Count(&teacherCount)
		assert.Equal(t, int64(param.Users), userCount)
		assert.Equal(t, int64(param.Courses), courseCount)
		assert.Equal(t, int64(param.Reviews), reviewCount)
		assert.Equal(t, int64(param.Teachers), teacherCount)
		courses := make([]po.CoursePO, 0)
		db.Model(&po.CoursePO{}).Find(&courses)
		for _, course := range courses[:5] {
			t.Logf("Course: %v", course)
		}
	})
}
