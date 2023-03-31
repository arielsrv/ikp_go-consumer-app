package container

import (
	"sync"

	"github.com/src/main/app/infrastructure/kvs"
	"github.com/src/main/app/model"
	"github.com/src/main/app/services"
)

var (
	consumerServiceOnce sync.Once
	consumerService     services.IConsumerService
)

func ProvideConsumerService() services.IConsumerService {
	consumerServiceOnce.Do(func() {
		consumerService = services.NewConsumerService(kvs.NewElasticCacheClient[model.AppStatusDTO](ProvideKVSClient()))
	})
	return consumerService
}
