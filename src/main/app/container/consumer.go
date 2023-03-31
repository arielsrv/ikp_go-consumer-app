package container

import (
	"runtime"
	"sync"

	"github.com/src/main/app/client"
	"github.com/src/main/app/config"
	"github.com/src/main/app/consumer"
	"github.com/src/main/app/infrastructure/queue"
	"github.com/src/main/app/log"
	"github.com/src/main/app/pusher"
)

var (
	topicHandlerOnce sync.Once
	topicConsumer    consumer.Consumer
)

func ProvideQueueConsumer() consumer.Consumer {
	topicHandlerOnce.Do(func() {
		rbPusher := config.ProvideRestClients().Get("target-client")
		pusherClient := client.NewHTTPPusherClient(rbPusher, config.String("pusher.target-endpoint"))
		httpPusher := pusher.NewHTTPPusher(pusherClient)

		queueClient, err := queue.NewClient(queue.Config{
			Name:     config.String("queues.orders.name"),
			URL:      config.String("queues.orders.url"),
			Parallel: config.TryInt("queues.orders.parallel", 10),
			Timeout:  config.TryInt("queues.orders.timeout", 1000),
		}, ProvideAWSConfig())

		if err != nil {
			log.Fatal(err)
		}

		topicConsumer = consumer.NewConsumer(consumer.Config{
			QueueService:     queueClient,
			Pusher:           httpPusher,
			Workers:          config.TryInt("consumers.orders.workers", runtime.NumCPU()-1),
			TaskResolverType: consumer.Async,
		}, ProvideConsumerService())
	})

	return topicConsumer
}
