package pusher

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/config"
	"github.com/src/main/app/config/env"
	"github.com/src/main/app/metrics"
	"github.com/src/main/app/rest"
	"log"
)

type Pusher interface {
	SendMessage(message *sqs.Message) error
}

type HttpPusher struct {
	httpClient rest.AppClient
}

type MessageDTO struct {
	Id      string `json:"MessageId,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewHttpPusher(httpClient rest.AppClient) *HttpPusher {
	return &HttpPusher{
		httpClient: httpClient,
	}
}

func (h HttpPusher) SendMessage(message *sqs.Message) error {
	var messageDTO MessageDTO
	err := json.Unmarshal([]byte(*message.Body), &messageDTO)
	if err != nil {
		log.Println(err)
		return err
	}

	requestBody := new(rest.RequestBody)
	requestBody.Id = messageDTO.Id
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
