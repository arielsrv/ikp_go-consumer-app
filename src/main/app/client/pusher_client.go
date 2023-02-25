package client

import (
	"net/http"
	"time"

	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/src/main/app/metrics"
	"github.com/src/main/app/server/errors"
)

type AppClient interface {
	PostMessage(body *RequestBody) error
}

type HTTPPusherClient struct {
	rb             rest.IRequestBuilder
	targetEndpoint string
}

func NewHTTPPusherClient(rb rest.IRequestBuilder, endpoint string) HTTPPusherClient {
	return HTTPPusherClient{
		rb:             rb,
		targetEndpoint: endpoint,
	}
}

func (c HTTPPusherClient) PostMessage(requestBody *RequestBody) error {
	startTime := time.Now()
	response := c.rb.Post(c.targetEndpoint, requestBody)
	elapsedTime := time.Since(startTime)

	metrics.Collector.RecordExecutionTime("consumers.pusher.clients.time", elapsedTime.Milliseconds())

	if response.Err != nil {
		return response.Err
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		metrics.Collector.IncrementCounter("consumers.pusher.clients.20x")
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		metrics.Collector.IncrementCounter("consumers.pusher.clients.40x")
	} else if response.StatusCode >= 500 {
		metrics.Collector.IncrementCounter("consumers.pusher.clients.50x")
	}

	if response.StatusCode != http.StatusOK {
		return errors.NewError(response.StatusCode, response.String())
	}

	return nil
}

type RequestBody struct {
	ID  string `json:"id,omitempty"`
	Msg string `json:"msg,omitempty"`
}
