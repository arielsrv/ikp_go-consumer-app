package pusher

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
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
		// @TODO: "h.metricCollector.IncrementCounter("consumer.app_name.pusher.errors")
		return err
	}
	// @TODO: "h.metricCollector.IncrementCounter("consumer.app_name.pusher.ok")
	return nil
}
