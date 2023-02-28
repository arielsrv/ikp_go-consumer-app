package container

import (
	"runtime"
	"sync"

	"github.com/src/main/app/client"
	"github.com/src/main/app/log"
	"github.com/src/main/app/pusher"
	"github.com/src/main/app/queue"

	"github.com/src/main/app/config"
	"github.com/src/main/app/consumer"
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
	topicConsumer       consumer.Consumer
)

func ProvideQueueConsumer() consumer.Consumer {
	consumerHandlerOnce.Do(func() {
		rbPusher := config.ProvideRestClients().Get("target-client")
		pusherClient := client.NewHTTPPusherClient(rbPusher, config.String("pusher.target-endpoint"))
		httpPusher := pusher.NewHTTPPusher(pusherClient)

		queueClient, err := queue.NewClient(queue.Config{
			QueueName: config.String("queues.orders.name"),
			Parallel:  config.TryInt("queues.orders.parallel", 10),
			Timeout:   config.TryInt("queues.orders.timeout", 1000),
		})

		if err != nil {
			log.Fatal(err)
		}

		topicConsumer = consumer.NewConsumer(consumer.Config{
			MessageClient: queueClient,
			Pusher:        httpPusher,
			Workers:       config.TryInt("consumers.orders.workers", runtime.NumCPU()-1)})
	})

	return topicConsumer
}
