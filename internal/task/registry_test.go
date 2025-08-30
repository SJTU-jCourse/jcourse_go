package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"jcourse_go/internal/task/base"
)

func TestInitTaskManager(t *testing.T) {
	InitTaskManager(base.RedisConfig{
		DSN:      "localhost:6379",
		Password: "",
	})
	time.Sleep(10 * time.Second)
	assert.NotNil(t, GlobalTaskManager)

	err := GlobalTaskManager.Shutdown()
	assert.Nil(t, err)
}
