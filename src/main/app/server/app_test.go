package server_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/src/main/app/routes"
	"github.com/src/main/app/server"
	"github.com/stretchr/testify/assert"
)

func TestApp_Start(t *testing.T) {
	app := server.New()
	routes.RegisterRoutes(app)

	port := 6789

	go func() {
		time.Sleep(500 * time.Millisecond)

		rb := rest.RequestBuilder{
			BaseURL: fmt.Sprintf("http://localhost:%d", port),
		}

		response := rb.Get("/ping")
		assert.NoError(t, response.Err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		err := app.Server.Shutdown()
		assert.NoError(t, err)
	}()

	assert.Nil(t, app.Start(fmt.Sprintf(":%d", port)))
}

func TestApp_Starter(t *testing.T) {
	app := server.New(server.Settings{
		Recovery:  true,
		Swagger:   true,
		RequestID: true,
		Logger:    true,
		Cors:      true,
		Metrics:   true,
	})
	routes.RegisterRoutes(app)

	listener, err := net.Listen("tcp", ":0")
	port := listener.Addr().(*net.TCPAddr).Port

	assert.NoError(t, err)

	go func() {
		time.Sleep(500 * time.Millisecond)

		rb := rest.RequestBuilder{
			BaseURL: fmt.Sprintf("http://localhost:%d", port),
		}

		response := rb.Get("/ping")
		assert.NoError(t, response.Err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		response = rb.Get("/metrics")
		assert.NoError(t, response.Err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		err = app.Server.Shutdown()
		assert.NoError(t, err)
	}()

	assert.Nil(t, app.Starter(listener))
}

func TestApp_StarterScope(t *testing.T) {
	t.Setenv("SCOPE", "prod")
	app := server.New(server.Settings{
		Swagger: true,
	})
	routes.RegisterRoutes(app)

	listener, err := net.Listen("tcp", ":0")
	port := listener.Addr().(*net.TCPAddr).Port

	assert.NoError(t, err)

	go func() {
		time.Sleep(500 * time.Millisecond)

		rb := rest.RequestBuilder{
			BaseURL: fmt.Sprintf("http://localhost:%d", port),
		}

		response := rb.Get("/ping")
		assert.NoError(t, response.Err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		err = app.Server.Shutdown()
		assert.NoError(t, err)
	}()

	assert.Nil(t, app.Starter(listener))
}
