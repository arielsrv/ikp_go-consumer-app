package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/src/main/app/log"
	"github.com/src/main/app/routes"
	"github.com/src/main/app/server"
)

func BenchmarkPingHandler_Ping(b *testing.B) {
	app := server.New()
	routes.RegisterRoutes(app)

	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(http.MethodGet, "/ping", nil)
		response, err := app.Server.Test(request)
		if err != nil || response.StatusCode != http.StatusOK {
			log.Infof("f[" + strconv.Itoa(i) + "] Status != OK (200)")
		}
	}
}
