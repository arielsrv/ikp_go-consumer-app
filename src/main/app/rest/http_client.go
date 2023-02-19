package rest

import (
	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/config"
	"github.com/src/main/app/metrics"
	"net/http"
	"time"
)

type AppClient interface {
	PostMessage(targetAppRequest *RequestBody) error
}

type HttpAppClient struct {
	rb      *rest.RequestBuilder
	baseURL string
}

func NewHttpAppClient(rb *rest.RequestBuilder) HttpAppClient {
	return HttpAppClient{
		rb:      rb,
		baseURL: config.String("target-app.endpoint"),
	}
}

func (c HttpAppClient) PostMessage(requestBody *RequestBody) error {
	startTime := time.Now()
	response := c.rb.Post(c.baseURL, requestBody)
	elapsed := time.Since(startTime)

	metrics.Collector.RecordExecutionTime("consumers.pusher.http.time", elapsed.Milliseconds())

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		metrics.Collector.IncrementCounter("consumers.pusher.http.20x")
	} else {
		if response.StatusCode >= 400 && response.StatusCode < 500 {
			metrics.Collector.IncrementCounter("consumers.pusher.http.40x")
		} else {
			if response.StatusCode >= 500 {
				metrics.Collector.IncrementCounter("consumers.pusher.http.50x")
			}
		}
	}

	if response.StatusCode != http.StatusOK {
		return fiber.NewError(response.StatusCode, response.String())
	}

	return nil
}

type RequestBody struct {
	Id  string `json:"id,omitempty"`
	Msg string `json:"msg,omitempty"`
}
