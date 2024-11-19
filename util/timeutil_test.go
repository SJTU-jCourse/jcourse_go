package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	t.Run("TestGetLocation", func(t *testing.T) {
		location := GetLocation()
		if location == nil {
			t.Errorf("GetLocation() error")
		}
		t.Logf("location: %v", location)
		assert.NotNil(t, location)

		mid := GetMidTime(time.Now())
		t.Logf("mid: %v", mid)
		start, end := GetDayTimeRange(time.Now())
		t.Logf("start: %v, end: %v", start, end)
	})
}
