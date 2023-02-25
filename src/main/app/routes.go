package app

import (
	"net/http"

	"github.com/src/main/app/handlers"
	"github.com/src/main/app/server"
)

func ProvideRoutes() []server.Route {
	var routes []server.Route

	routes = append(routes, server.
		Route{Verb: http.MethodGet, Path: "/ping", Action: server.Resolve[handlers.PingHandler]().Ping})

	return routes
}
