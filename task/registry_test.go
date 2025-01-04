package task

import (
	"jcourse_go/task/base"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
