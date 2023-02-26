package app

import (
	"context"
	"fmt"

	"github.com/src/main/app/container"
	"github.com/src/main/app/routes"

	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/log"

	"github.com/src/main/app/server"
)

func Run() error {
	app := server.New(server.Settings{
		Metrics: true,
		Swagger: true,
	})

	routes.RegisterRoutes(app)

	if !env.IsProd() {
		topicConsumer := container.ProvideQueueConsumer()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go topicConsumer.Start(ctx)
	}

	host := config.String("HOST")
	if env.IsEmpty(host) && !env.IsLocal() {
		host = "0.0.0.0"
	}

	port := config.String("PORT")
	if env.IsEmpty(port) {
		port = "8080"
	}

	address := fmt.Sprintf("%s:%s", host, port)
	log.Infof("Listening on local address %s", address)

	log.Infof(fmt.Sprintf("Open %s/ping in the browser", config.String("public")))

	return app.Start(address)
}
