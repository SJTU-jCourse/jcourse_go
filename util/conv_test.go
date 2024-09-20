package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringsToFloat64s(t *testing.T) {
	src := []string{"1", "2", "3", "4", "5"}
	dst := StringsToFloat64s(src)
	assert.Equal(t, len(src), len(dst))
	for i := range src {
		assert.Equal(t, float64(i+1), dst[i])
	}
}
