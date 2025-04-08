package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/model/types"
	"jcourse_go/repository"
)

func TestSyncRating(t *testing.T) {
	repository.SetupTestEnv()
	ctx := context.Background()
	t.Run("sync course", func(t *testing.T) {
		c := repository.Q.CoursePO

		coursePO, err := c.WithContext(ctx).Where(c.ID.Eq(1)).Take()
		assert.Nil(t, err)
		assert.Equal(t, 0, int(coursePO.RatingCount))

		err = SyncRating(ctx, types.RelatedTypeCourse, coursePO.ID)
		assert.Nil(t, err)

		coursePO, err = c.WithContext(ctx).Where(c.ID.Eq(1)).Take()
		assert.Nil(t, err)
		assert.Equal(t, 1, int(coursePO.RatingCount))
		assert.Equal(t, 5.0, coursePO.RatingAvg)
	})

	t.Run("sync teacher", func(t *testing.T) {
		te := repository.Q.TeacherPO

		teacherPO, err := te.WithContext(ctx).Where(te.ID.Eq(2)).Take()
		assert.Nil(t, err)
		assert.Equal(t, 0, int(teacherPO.RatingCount))

		err = SyncRating(ctx, types.RelatedTypeTeacher, teacherPO.ID)
		assert.Nil(t, err)

		teacherPO, err = te.WithContext(ctx).Where(te.ID.Eq(2)).Take()
		assert.Nil(t, err)
		assert.Equal(t, 1, int(teacherPO.RatingCount))
		assert.Equal(t, 5.0, teacherPO.RatingAvg)
	})
	repository.TearDownTestEnv()
}
