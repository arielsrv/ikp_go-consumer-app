package routes

import (
	"net/http"

	"github.com/src/main/app/container"
	"github.com/src/main/app/server"
)

func RegisterRoutes(app *server.App) {
	app.Route(http.MethodGet, "/ping", container.ProvidePingHandler().Ping)
	app.Route(http.MethodGet, "/consumer/status", container.ProvideConsumerHandler().GetStatus)
	app.Route(http.MethodPut, "/consumer/start", container.ProvideConsumerHandler().Start)
	app.Route(http.MethodPut, "/consumer/stop", container.ProvideConsumerHandler().Stop)
}
