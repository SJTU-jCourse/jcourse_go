package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/dal"
)

func TestTrainingPlanQuery_GetTrainingPlan(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewTrainingPlanQuery(db)
	t.Run("has", func(t *testing.T) {
		tps, err := query.GetTrainingPlan(ctx)
		assert.Nil(t, err)
		assert.Len(t, tps, 1)
		tp := tps[0]
		assert.Equal(t, "计算机科学与技术", tp.Major)
	})
	t.Run("none", func(t *testing.T) {
		tps, err := query.GetTrainingPlan(ctx, WithID(20))
		assert.Nil(t, err)
		assert.Len(t, tps, 0)
	})
}

func TestTrainingPlanQuery_GetTrainingPlanCount(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewTrainingPlanQuery(db)
	t.Run("has", func(t *testing.T) {
		count, err := query.GetTrainingPlanCount(ctx)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	})
	t.Run("none", func(t *testing.T) {
		count, err := query.GetTrainingPlanCount(ctx, WithID(20))
		assert.Nil(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestTrainingPlanCourseQuery_GetTrainingPlanCourseList(t *testing.T) {
	ctx := context.Background()
	db := dal.GetDBClient()
	query := NewTrainingPlanCourseQuery(db)
	t.Run("has", func(t *testing.T) {
		courses, err := query.GetTrainingPlanCourseList(ctx, WithTrainingPlanID(1))
		assert.Nil(t, err)
		assert.Len(t, courses, 3)
	})
	t.Run("none", func(t *testing.T) {
		courses, err := query.GetTrainingPlanCourseList(ctx, WithID(20))
		assert.Nil(t, err)
		assert.Len(t, courses, 0)
	})
}
