package pusher

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/metrics"
	"github.com/src/main/app/rest"
)

type Pusher interface {
	SendMessage(message *sqs.Message) error
}

type HttpPusher struct {
	httpClient rest.IHttpClient
}

type MessageDTO struct {
	Message string
}

func NewHttpPusher(httpClient rest.IHttpClient) *HttpPusher {
	return &HttpPusher{
		httpClient: httpClient,
	}
}

func (h HttpPusher) SendMessage(message *sqs.Message) error {
	var messageDTO MessageDTO
	err := json.Unmarshal([]byte(*message.Body), &messageDTO)
	if err != nil {
		return nil
	}

	requestBody := new(rest.RequestBody)
	requestBody.Msg = messageDTO.Message
	err = h.httpClient.PostMessage(requestBody)

	if err != nil {
		metrics.Collector.IncrementCounter("consumers.pusher.errors",
			fmt.Sprintf("name: %s", config.String("app.name")),
			fmt.Sprintf("scope: %s", env.GetScope()))
		return err
	}

	metrics.Collector.IncrementCounter("consumers.pusher.success",
		fmt.Sprintf("name: %s", config.String("app.name")),
		fmt.Sprintf("scope: %s", env.GetScope()))

	return nil
}
