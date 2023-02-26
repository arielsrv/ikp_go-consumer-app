package routes

import (
	"net/http"

	"github.com/src/main/app/container"
	"github.com/src/main/app/server"
)

func RegisterRoutes(app *server.App) {
	app.Route(http.MethodGet, "/ping", container.ProvidePingHandler().Ping)
}
