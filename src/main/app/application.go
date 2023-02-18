package app

import (
	"context"
	"fmt"
	"github.com/src/main/app/clients"
	"github.com/src/main/app/consumer"
	"github.com/src/main/app/infrastructure"
	"log"
	"net/http"
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
		Swagger:   true,
		RequestID: true,
		Logger:    true,
		NewRelic:  true,
	})

	pingService := services.NewPingService()
	pingHandler := handlers.NewPingHandler(pingService)

	targetAppClient := clients.NewClient(restClients.Get("target-app"))
	log.Println(targetAppClient)

	server.RegisterHandler(pingHandler)
	server.Register(http.MethodGet, "/ping", server.Resolve[handlers.PingHandler]().Ping)

	consume()

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

	// Instantiate client.
	client := infrastructure.NewSQS(session, time.Second*5)

	// Instantiate consumer and start consuming.
	consumer.NewConsumer(client, consumer.Config{
		Type:      consumer.AsyncConsumer,
		QueueURL:  config.String("consumers.users.queue-url"),
		MaxWorker: config.TryInt("consumers.users.max-workers", 2),
		MaxMsg:    config.TryInt("consumers.users.max-messages", 10),
	}).Start(ctx)
}
