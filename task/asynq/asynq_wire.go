package impl

import (
	"jcourse_go/task"
	"jcourse_go/util"
)

// TODO @huangjunqing wire
// import (
// 	"github.com/google/wire"
// )

// // ProviderSet is a Wire provider set that provides IAsyncTaskManager.
// var ProviderSet = wire.NewSet(
// 	NewAsynqTaskManager,
// 	wire.Bind(new(IAsyncTaskManager), new(*AsynqTaskManager)),
// )

// InitializeTaskManager initializes the IAsyncTaskManager interface.
func InitializeTaskManager() (task.IAsyncTaskManager, error) {
	manager := NewAsynqTaskManager(redisConfig{
		Host:     util.GetRedisHost(),
		Port:     util.GetRedisPort(),
		Password: util.GetRedisPassword(),
	})
	err := manager.StartServer()
	if err != nil {
		return nil, err
	}
	return manager, nil
}
