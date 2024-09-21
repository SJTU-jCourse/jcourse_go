package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/dal"
)

func TestTeacherQuery_GetTeacher(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewTeacherQuery(db)
	teachers, err := query.GetTeacher(ctx, WithID(1))
	assert.Nil(t, err)
	assert.Len(t, teachers, 1)

	teacher := teachers[0]
	assert.Equal(t, "10001", teacher.Code)
	assert.Equal(t, "高女士", teacher.Name)
}

func TestTeacherQuery_GetTeacherCount(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewTeacherQuery(db)
	t.Run("all", func(t *testing.T) {
		count, err := query.GetTeacherCount(ctx)
		assert.Nil(t, err)
		assert.Equal(t, int64(4), count)
	})
	t.Run("filter count", func(t *testing.T) {
		count, err := query.GetTeacherCount(ctx, WithCode("10001"))
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	})
	t.Run("none", func(t *testing.T) {
		count, err := query.GetTeacherCount(ctx, WithCode("1001"))
		assert.Nil(t, err)
		assert.Equal(t, int64(0), count)
	})
}
