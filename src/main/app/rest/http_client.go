package rest

import (
	"fmt"
	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/metrics"
	"net/http"
	"time"
)

type IHttpClient interface {
	PostMessage(targetAppRequest *RequestBody) error
}

type Client struct {
	rb      *rest.RequestBuilder
	baseURL string
}

func NewClient(rb *rest.RequestBuilder) Client {
	return Client{
		rb:      rb,
		baseURL: config.String("target-app.endpoint"),
	}
}

func (c Client) PostMessage(requestBody *RequestBody) error {
	startTime := time.Now()
	response := c.rb.Post(c.baseURL, requestBody)
	elapsed := time.Since(startTime)

	metrics.Collector.RecordExecutionTime("consumers.pusher.http.time",
		elapsed.Milliseconds(), "name: %s", config.String("app.name"))

	if response.Err != nil {
		return response.Err
	}

	if response.StatusCode != http.StatusOK {
		return fiber.NewError(response.StatusCode, response.String())
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		metrics.Collector.IncrementCounter("consumers.pusher.http.20x",
			fmt.Sprintf("name: %s", config.String("app.name")),
			fmt.Sprintf("scope: %s", env.GetScope()))
	} else {
		if response.StatusCode >= 400 && response.StatusCode < 500 {
			metrics.Collector.IncrementCounter("consumers.pusher.http.40x",
				fmt.Sprintf("name: %s", config.String("app.name")),
				fmt.Sprintf("scope: %s", env.GetScope()))
		} else {
			if response.StatusCode >= 500 {
				metrics.Collector.IncrementCounter("consumers.pusher.http.50x",
					fmt.Sprintf("name: %s", config.String("app.name")),
					fmt.Sprintf("scope: %s", env.GetScope()))
			}
		}
	}

	return nil
}

type RequestBody struct {
	Id  string `json:"id,omitempty"`
	Msg string `json:"msg,omitempty"`
}
