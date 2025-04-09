package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/model/dto"
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

func TestCreateRating(t *testing.T) {
	repository.SetupTestEnv()
	defer repository.TearDownTestEnv()

	tests := []struct {
		name          string
		userID        int64
		ratingDTO     dto.RatingDTO
		expectedError bool
		verify        func(*testing.T)
	}{
		{
			name:   "create new course rating",
			userID: 1,
			ratingDTO: dto.RatingDTO{
				RelatedType: string(types.RelatedTypeCourse),
				RelatedID:   1,
				Rating:      5,
			},
			expectedError: false,
			verify: func(t *testing.T) {
				coursePO, err := repository.Q.CoursePO.WithContext(context.Background()).Where(repository.Q.CoursePO.ID.Eq(1)).Take()
				assert.Nil(t, err)
				assert.Equal(t, 1, int(coursePO.RatingCount))
				assert.Equal(t, 5.0, coursePO.RatingAvg)
			},
		},
		{
			name:   "update existing course rating",
			userID: 1,
			ratingDTO: dto.RatingDTO{
				RelatedType: string(types.RelatedTypeCourse),
				RelatedID:   1,
				Rating:      3,
			},
			expectedError: false,
			verify: func(t *testing.T) {
				coursePO, err := repository.Q.CoursePO.WithContext(context.Background()).Where(repository.Q.CoursePO.ID.Eq(1)).Take()
				assert.Nil(t, err)
				assert.Equal(t, 1, int(coursePO.RatingCount))
				assert.Equal(t, 3.0, coursePO.RatingAvg)
			},
		},
		{
			name:   "invalid type",
			userID: 1,
			ratingDTO: dto.RatingDTO{
				RelatedType: "invalid",
				RelatedID:   1,
				Rating:      4,
			},
			expectedError: true,
			verify:        func(t *testing.T) {},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CreateRating(context.Background(), test.userID, test.ratingDTO)
			if test.expectedError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			if test.verify != nil {
				test.verify(t)
			}
		})
	}
}
