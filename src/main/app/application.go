package app

import (
	"context"
	"fmt"
	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/consumer"
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/log"
	"github.com/src/main/app/pusher"
	"github.com/src/main/app/queue"
	"github.com/src/main/app/rest"
	"github.com/src/main/app/server"
	"github.com/src/main/app/services"
	"net/http"
	"runtime"
)

var restClients = config.ProvideRestClients()

func Run() error {
	app := server.New(server.Config{
		Recovery:  true,
		RequestID: true,
		Logger:    true,
		Metrics:   true,
	})

	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)
	server.RegisterHandler(pingHandler)
	server.Register(http.MethodGet, "/ping", server.Resolve[handlers.PingHandler]().Ping)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requestBuilder := restClients.Get("target-client")
	httpClient := rest.NewHttpAppClient(requestBuilder, config.String("pusher.target-endpoint"))
	httpPusher := pusher.NewHttpPusher(httpClient)

	queueClient, err := queue.
		NewClient(
			config.String("queues.users.name"),
			config.TryInt("queues.users.parallel", 10),
			config.TryInt("queues.users.timeout", 1000))

	if err != nil {
		log.Fatal(err)
	}

	go consumer.
		NewConsumer(
			queueClient,
			httpPusher,
			config.TryInt("consumers.users.workers", runtime.
				NumCPU()-1)).
		Start(ctx)

	host := config.String("HOST")
	if env.IsEmpty(host) && !env.IsLocal() {
		host = "0.0.0.0"
	} else {
		host = "127.0.0.1"
	}

	port := config.String("PORT")
	if env.IsEmpty(port) {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)

	log.Infof("Listening on port %s", port)
	log.Infof("Open http://%s:%s/ping in the browser", host, port)

	return app.Start(address)
}
