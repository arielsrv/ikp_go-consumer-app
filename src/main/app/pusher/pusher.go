package pusher

import (
	"encoding/json"

	"github.com/src/main/app/client"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/log"
	"github.com/src/main/app/metrics"
)

type Pusher interface {
	SendMessage(message *sqs.Message) error
}

type HTTPPusher struct {
	httpClient client.AppClient
}

type MessageDTO struct {
	ID      string `json:"MessageId,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewHTTPPusher(httpClient client.AppClient) *HTTPPusher {
	return &HTTPPusher{
		httpClient: httpClient,
	}
}

func (h HTTPPusher) SendMessage(message *sqs.Message) error {
	var messageDTO MessageDTO
	err := json.Unmarshal([]byte(*message.Body), &messageDTO)
	if err != nil {
		log.Error(err)
		return err
	}

	requestBody := new(client.RequestBody)
	requestBody.ID = messageDTO.ID
	requestBody.Msg = messageDTO.Message

	log.Infof("message - id: %s, body: %s", requestBody.ID, requestBody.Msg)

	err = h.httpClient.PostMessage(requestBody)

	if err != nil {
		metrics.Collector.IncrementCounter("consumers.pusher.errors")
		return err
	}

	metrics.Collector.IncrementCounter("consumers.pusher.success")

	return nil
}
