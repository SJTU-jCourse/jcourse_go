package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/model/model"
)

func TestRemoveReviewUserInfo(t *testing.T) {
	hides := []bool{true, false, true, false}
	expected := []int64{0, 1, 1, 1}
	user := model.UserMinimal{ID: 1}
	reviews := []model.Review{{
		User:        user,
		IsAnonymous: true,
	}, {
		User:        user,
		IsAnonymous: true,
	}, {
		User:        user,
		IsAnonymous: false,
	}, {
		User:        user,
		IsAnonymous: false,
	}}
	for i := range reviews {
		r := reviews[i]
		hide := hides[i]
		RemoveReviewUserInfo(&r, hide)
		assert.Equal(t, expected[i], r.User.ID)
	}
}

func TestRemoveReviewsUserInfo(t *testing.T) {
	reviews := []model.Review{{
		User:        model.UserMinimal{ID: 1},
		IsAnonymous: true,
	}}
	RemoveReviewsUserInfo(reviews, true)
	assert.Equal(t, int64(0), reviews[0].User.ID)
}
