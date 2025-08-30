package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRatingInfo_Calc(t *testing.T) {
	rating := &RatingInfo{
		RatingDist: []RatingInfoDistItem{
			{Count: 1, Rating: 1},
			{Count: 2, Rating: 2},
			{Count: 5, Rating: 5},
		},
	}
	rating.Calc()
	assert.Equal(t, int64(8), rating.Count)
	assert.Equal(t, 3.75, rating.Average)
}
