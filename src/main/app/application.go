package app

import (
	"context"
	"fmt"
	"github.com/src/main/app/clients"
	"github.com/src/main/app/consumer"
	"github.com/src/main/app/infrastructure"
	"github.com/src/main/app/pusher"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/handlers"
	"github.com/src/main/app/server"
	"github.com/src/main/app/services"
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

	go consume()

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

func consume() {
	// Create a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a session instance.
	session, err := infrastructure.New(infrastructure.Config{
		Address: config.String("aws.url"),
		Region:  config.String("aws.region"),
		Profile: config.String("aws.profile"),
		ID:      config.String("aws.id"),
		Secret:  config.String("aws.secret"),
	})
	if err != nil {
		log.Fatalln(err)
	}

	messageClient := infrastructure.NewSQS(session, time.Second*5)
	httpClient := clients.NewClient(restClients.Get("target-app"))
	httpPusher := pusher.NewHttpPusher(httpClient)

	// Instantiate consumer and start consuming.
	consumer.NewConsumer(messageClient, httpPusher, consumer.Config{
		QueueURL: config.String("consumers.users.queue-url"),
		Workers:  config.TryInt("consumers.users.workers", runtime.NumCPU()-1),
		MaxMsg:   config.TryInt("consumers.users.workers.messages", 10),
	}).Start(ctx)
}
