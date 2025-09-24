package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/internal/domain/course"
	"jcourse_go/internal/domain/user"
)

func TestRemoveReviewUserInfo(t *testing.T) {
	hides := []bool{true, false, true, false}
	expected := []int64{0, 1, 1, 1}
	user := user.UserMinimal{ID: 1}
	reviews := []course.Review{{
		User:     user,
		IsPublic: true,
	}, {
		User:     user,
		IsPublic: true,
	}, {
		User:     user,
		IsPublic: false,
	}, {
		User:     user,
		IsPublic: false,
	}}
	for i := range reviews {
		r := reviews[i]
		hide := hides[i]
		RemoveReviewUserInfo(&r, 0, hide)
		assert.Equal(t, expected[i], r.User.ID)
	}
}

func TestRemoveReviewsUserInfo(t *testing.T) {
	reviews := []course.Review{{
		User:     user.UserMinimal{ID: 1},
		IsPublic: true,
	}}
	RemoveReviewsUserInfo(reviews, 0, true)
	assert.Equal(t, int64(0), reviews[0].User.ID)
}
