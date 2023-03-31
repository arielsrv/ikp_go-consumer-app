package container

import (
	"sync"

	"github.com/src/main/app/handlers"
	"github.com/src/main/app/services"
)

var (
	pingHandlerOnce sync.Once
	pingHandler     *handlers.PingHandler
)

func ProvidePingHandler() *handlers.PingHandler {
	pingHandlerOnce.Do(func() {
		pingService := services.NewPingService()
		pingHandler = handlers.NewPingHandler(pingService)
	})

	return pingHandler
}

var (
	consumerHandlerOnce sync.Once
	consumerHandler     *handlers.ConsumerHandler
)

func ProvideConsumerHandler() *handlers.ConsumerHandler {
	consumerHandlerOnce.Do(func() {
		consumerHandler = handlers.NewConsumerHandler(ProvideConsumerService())
	})
	return consumerHandler
}
