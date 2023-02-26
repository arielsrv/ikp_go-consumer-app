package server

import (
	"fmt"
	"net/http"
	"os"

	properties "github.com/src/main/app/config"
	"github.com/src/main/app/config/env"

	"github.com/gofiber/fiber/v2"

	"github.com/src/main/app/log"

	"github.com/arielsrv/nrfiber"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/src/main/app/server/errors"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
)

var routes []Route

type Settings struct {
	Recovery  bool
	Swagger   bool
	RequestID bool
	Logger    bool
	Cors      bool
	NewRelic  bool
	Metrics   bool
}

type App struct {
	*fiber.App
	config Settings
}

type Route struct {
	Verb   string
	Path   string
	Action func(ctx *fiber.Ctx) error
}

func (app *App) Start(addr string) error {
	for _, route := range routes {
		app.Add(route.Verb, route.Path, route.Action)
	}
	return app.Listen(addr)
}

func (app *App) Route(method, path string, handlers ...fiber.Handler) {
	app.Add(method, path, handlers...)
}

func New(config ...Settings) *App {
	app := &App{
		App: fiber.New(fiber.Config{
			DisableStartupMessage: true,
			ErrorHandler:          errors.ErrorHandler,
		}),
		config: Settings{
			Recovery:  true,
			Swagger:   false,
			RequestID: false,
			Logger:    false,
			Cors:      false,
			NewRelic:  false,
			Metrics:   false,
		},
	}

	if len(config) > 0 {
		app.config = config[0]
	}

	if app.config.Recovery {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	if app.config.RequestID {
		app.Use(requestid.New())
	}

	if app.config.Logger {
		app.Use(logger.New(logger.Config{
			Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
		}))
	}

	if app.config.Cors {
		app.Use(cors.New())
	}

	if app.config.Swagger {
		if !env.IsLocal() {
			app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
				URL: fmt.Sprintf("%s/swagger/doc.json",
					properties.String("public")),
			}))
		} else {
			app.Add(http.MethodGet, "/swagger/*", swagger.HandlerDefault)
		}
		log.Info("Swagger enabled")
	}

	if app.config.NewRelic && !env.IsLocal() {
		newRelicLicense := properties.String("NEW_RELIC_LICENSE_KEY")
		if !env.IsEmpty(newRelicLicense) {
			nrApp, err := newrelic.NewApplication(
				newrelic.ConfigAppName("app.name"),
				newrelic.ConfigLicense(newRelicLicense),
				newrelic.ConfigDebugLogger(os.Stdout),
			)
			if err != nil {
				log.Fatalf("failed to load newrelic config %s", err)
			}

			app.Use(nrfiber.New(nrfiber.Config{
				NewRelicApp: nrApp,
			}))
		} else {
			log.Info("newrelic disabled, newrelic license key not found")
		}
	}

	if app.config.Metrics {
		prometheus := fiberprometheus.New(properties.String("app.name"))
		prometheus.RegisterAt(app.App, "/metrics")
		app.Use(prometheus.Middleware)
	}

	return app
}
