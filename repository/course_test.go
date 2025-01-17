package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/dal"
)

func TestBaseCourseQuery_GetBaseCourse(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewBaseCourseQuery(db)
	courses, err := query.GetBaseCourse(ctx, WithCode("CS1500"))
	assert.Nil(t, err)
	assert.Len(t, courses, 1)
	assert.Equal(t, "CS1500", courses[0].Code)
	assert.Equal(t, "计算机科学导论", courses[0].Name)
	assert.Equal(t, float64(4), courses[0].Credit)
	assert.NotEqual(t, 0, courses[0].ID)
}

func TestBaseCourseQuery_GetBaseCourseCount(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewBaseCourseQuery(db)
	t.Run("all", func(t *testing.T) {
		count, err := query.GetBaseCourseCount(ctx)
		assert.Nil(t, err)
		assert.Equal(t, int64(3), count)
	})
	t.Run("filter count", func(t *testing.T) {
		count, err := query.GetBaseCourseCount(ctx, WithCode("CS1500"))
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	})
}

func TestBaseCourseQuery_GetCourseByIDs(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewBaseCourseQuery(db)
	t.Run("all", func(t *testing.T) {
		courseMap, err := query.GetBaseCoursesByIDs(ctx, []int64{1, 2})
		assert.Nil(t, err)
		assert.Len(t, courseMap, 2)
		assert.Equal(t, uint(1), courseMap[1].ID)
		assert.Equal(t, "MARX1001", courseMap[1].Code)
		assert.Equal(t, uint(2), courseMap[2].ID)
		assert.Equal(t, "CS1500", courseMap[2].Code)
	})
	t.Run("not found", func(t *testing.T) {
		courseMap, err := query.GetBaseCoursesByIDs(ctx, []int64{10, 20})
		assert.Nil(t, err)
		assert.Len(t, courseMap, 0)
	})
}

func TestCourseQuery_GetCourse(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewCourseQuery(db)
	t.Run("all", func(t *testing.T) {
		courses, err := query.GetCourse(ctx)
		assert.Nil(t, err)
		assert.Len(t, courses, 4)
	})
	t.Run("filter course", func(t *testing.T) {
		courses, err := query.GetCourse(ctx, WithCode("CS1500"))
		assert.Nil(t, err)
		assert.Len(t, courses, 1)
		assert.Equal(t, "CS1500", courses[0].Code)
	})
	t.Run("none", func(t *testing.T) {
		courses, err := query.GetCourse(ctx, WithCode("CS3500"))
		assert.Nil(t, err)
		assert.Len(t, courses, 0)
	})
}

func TestCourseQuery_GetCourseCount(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewCourseQuery(db)
	t.Run("all", func(t *testing.T) {
		count, err := query.GetCourseCount(ctx)
		assert.Nil(t, err)
		assert.Equal(t, int64(4), count)
	})
	t.Run("filter course", func(t *testing.T) {
		count, err := query.GetCourseCount(ctx, WithCode("CS1500"))
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	})
	t.Run("none", func(t *testing.T) {
		count, err := query.GetCourseCount(ctx, WithCode("CS3500"))
		assert.Nil(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestCourseQuery_GetCourseByIDs(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewCourseQuery(db)
	courseMap, err := query.GetCourseByIDs(ctx, []int64{1, 2})
	assert.Nil(t, err)
	assert.Len(t, courseMap, 2)
	assert.Equal(t, uint(1), courseMap[1].ID)
	assert.Equal(t, "MARX1001", courseMap[1].Code)
	assert.Equal(t, uint(2), courseMap[2].ID)
	assert.Equal(t, "MARX1001", courseMap[2].Code)
}

func TestCourseQuery_GetCourseCategories(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewCourseQuery(db)

	t.Run("no category", func(t *testing.T) {
		categories, err := query.GetCourseCategories(ctx, []int64{3})
		assert.Nil(t, err)
		assert.Len(t, categories, 1)
		assert.Len(t, categories[3], 0)
	})

	t.Run("has category", func(t *testing.T) {
		categories, err := query.GetCourseCategories(ctx, []int64{2})
		assert.Nil(t, err)
		assert.Len(t, categories, 1)
		category := categories[2]
		assert.Len(t, category, 2)
		assert.Contains(t, category, "通识")
		assert.Contains(t, category, "必修")
	})
}
