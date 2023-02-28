package client

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/src/main/app/log"
	"github.com/src/main/app/server"

	"github.com/src/main/app/metrics"

	"github.com/arielsrv/ikp_go-restclient/rest"
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

	metrics.Collector.RecordExecutionTime(metrics.PusherHTTPTime, elapsedTime.Milliseconds())

	if response.Err != nil {
		var err net.Error
		if ok := errors.As(response.Err, &err); ok && err.Timeout() {
			log.Warnf("pusher timeout, discuss cap theorem, possible inconsistency ensure handle duplicates from target app, MessageId: %s", requestBody.ID)
			metrics.Collector.IncrementCounter(metrics.PusherHTTPTimeout)
		}
		return response.Err
	}

	switch {
	case response.StatusCode >= 200 && response.StatusCode < 300:
		metrics.Collector.IncrementCounter(metrics.PusherStatusOK)
	case response.StatusCode >= 400 && response.StatusCode < 500:
		metrics.Collector.IncrementCounter(metrics.PusherStatus40x)
	case response.StatusCode >= http.StatusInternalServerError:
		metrics.Collector.IncrementCounter(metrics.PusherStatus50x)
	}

	if response.StatusCode != http.StatusOK {
		return server.NewError(response.StatusCode, response.String())
	}

	return nil
}

type RequestBody struct {
	ID        string `json:"id,omitempty"`
	Msg       string `json:"msg,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}
