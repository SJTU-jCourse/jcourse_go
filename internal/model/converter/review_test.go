package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	model2 "jcourse_go/internal/model/model"
)

func TestRemoveReviewUserInfo(t *testing.T) {
	hides := []bool{true, false, true, false}
	expected := []int64{0, 1, 1, 1}
	user := model2.UserMinimal{ID: 1}
	reviews := []model2.Review{{
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
		RemoveReviewUserInfo(&r, 0, hide)
		assert.Equal(t, expected[i], r.User.ID)
	}
}

func TestRemoveReviewsUserInfo(t *testing.T) {
	reviews := []model2.Review{{
		User:        model2.UserMinimal{ID: 1},
		IsAnonymous: true,
	}}
	RemoveReviewsUserInfo(reviews, 0, true)
	assert.Equal(t, int64(0), reviews[0].User.ID)
}
