package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/dal"
)

func setup() {
	ctx := context.Background()
	dal.InitTestMemDBClient()
	db := dal.GetDBClient()
	_ = Migrate(db)
	_ = CreateTestEnv(ctx, db)
}

func tearDown() {
	db, _ := dal.GetDBClient().DB()
	_ = db.Close()
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tearDown()
}

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

func TestCourseQuery_GetCourseByIDs(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewBaseCourseQuery(db)
	courseMap, err := query.GetBaseCoursesByIDs(ctx, []int64{1, 2})
	assert.Nil(t, err)
	assert.Len(t, courseMap, 2)
	assert.Equal(t, uint(1), courseMap[1].ID)
	assert.Equal(t, "MARX1001", courseMap[1].Code)
	assert.Equal(t, uint(2), courseMap[2].ID)
	assert.Equal(t, "CS1500", courseMap[2].Code)
}
