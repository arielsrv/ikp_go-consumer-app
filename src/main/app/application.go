package app

import (
	"context"
	"fmt"
	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/consumer"
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/pusher"
	"github.com/src/main/app/queue"
	"github.com/src/main/app/rest"
	"github.com/src/main/app/server"
	"github.com/src/main/app/services"
	"log"
	"net/http"
)

var restClients = config.ProvideRestClients()

func Run() error {
	app := server.New(server.Config{
		Recovery:  true,
		RequestID: true,
		Logger:    true,
	})

	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	server.RegisterHandler(pingHandler)
	server.Register(http.MethodGet, "/ping", server.Resolve[handlers.PingHandler]().Ping)

	// Create a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient := rest.NewHttpAppClient(restClients.Get("target-app"))
	httpPusher := pusher.NewHttpPusher(httpClient)

	queueClient := queue.NewClient(config.String("consumers.users.queue-url"))
	consumer.NewConsumer(queueClient, httpPusher).Start(ctx)

	host := config.String("HOST")
	if env.IsEmpty(host) && !env.IsDev() {
		host = "0.0.0.0"
	} else {
		host = "127.0.0.1"
	}

	port := config.String("PORT")
	if env.IsEmpty(port) {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://%s:%s/ping in the browser", host, port)

	return app.Start(address)
}
