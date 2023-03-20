package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	properties "github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/log"
)

type Settings struct {
	Recovery  bool
	Swagger   bool
	RequestID bool
	Logger    bool
	Cors      bool
	Metrics   bool
}

type App struct {
	Server *fiber.App
	config Settings
}

type Route struct {
	Verb   string
	Path   string
	Action func(ctx *fiber.Ctx) error
}

func (app *App) Start(addr string) error {
	return app.Server.Listen(addr)
}

func (app *App) Starter(listener net.Listener) error {
	return app.Server.Listener(listener)
}

func (app *App) Route(method, path string, handlers ...fiber.Handler) {
	app.Server.Add(method, path, handlers...)
}

func New(config ...Settings) *App {
	app := &App{
		Server: fiber.New(fiber.Config{
			DisableStartupMessage: true,
			ErrorHandler:          ErrorHandler,
		}),
		config: Settings{
			Recovery:  true,
			Swagger:   false,
			RequestID: false,
			Logger:    false,
			Cors:      false,
			Metrics:   false,
		},
	}

	if len(config) > 0 {
		app.config = config[0]
	}

	if app.config.Recovery {
		app.Server.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	if app.config.RequestID {
		app.Server.Use(requestid.New())
	}

	if app.config.Logger {
		app.Server.Use(logger.New(logger.Config{
			Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
			Output: log.GetWriter(),
		}))
	}

	if app.config.Cors {
		app.Server.Use(cors.New())
	}

	if app.config.Swagger {
		if !env.IsLocal() {
			app.Server.Get("/swagger/*", swagger.New(swagger.Config{ // custom
				URL: fmt.Sprintf("%s/swagger/doc.json", properties.String("public")),
			}))
		} else {
			app.Server.Add(http.MethodGet, "/swagger/*", swagger.HandlerDefault)
		}
		log.Info("swagger enabled")
	}

	if app.config.Metrics {
		prometheus := fiberprometheus.New(properties.String("app.name"))
		prometheus.RegisterAt(app.Server, "/metrics")
		app.Server.Use(prometheus.Middleware)
	}

	return app
}
