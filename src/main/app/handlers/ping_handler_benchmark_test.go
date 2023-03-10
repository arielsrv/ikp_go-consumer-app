package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/src/main/app/log"

	"github.com/src/main/app/server"

	"github.com/src/main/app/handlers"
)

func BenchmarkPingHandler_Ping(b *testing.B) {
	pingService := new(MockPingService)
	pingHandler := handlers.NewPingHandler(pingService)
	app := server.New(server.Settings{
		Logger: false,
	})
	app.Server.Add(http.MethodGet, "/ping", pingHandler.Ping)

	pingService.On("Ping").Return("pong")

	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(http.MethodGet, "/ping", nil)
		response, err := app.Server.Test(request)
		if err != nil || response.StatusCode != http.StatusOK {
			log.Infof("f[" + strconv.Itoa(i) + "] Status != OK (200)")
		}
	}
}
