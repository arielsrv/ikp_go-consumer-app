package app

import (
	"runtime"

	"github.com/src/main/app/client"
	"github.com/src/main/app/config"
	"github.com/src/main/app/consumer"
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/log"
	"github.com/src/main/app/pusher"
	"github.com/src/main/app/queue"
	"github.com/src/main/app/services"
)

var restClients = config.ProvideRestClients()

func ProvidePingHandler() *handlers.PingHandler {
	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	return pingHandler
}
func ProvideQueueConsumer() consumer.Consumer {
	rbPusher := restClients.Get("target-client")
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

	topicConsumer := consumer.NewConsumer(consumer.Config{
		MessageClient: queueClient,
		Pusher:        httpPusher,
		Workers:       config.TryInt("consumers.orders.workers", runtime.NumCPU()-1)})

	return topicConsumer
}
