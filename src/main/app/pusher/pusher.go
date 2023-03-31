package pusher

import (
	"encoding/json"

	"github.com/src/main/app/client"
	"github.com/src/main/app/infrastructure/queue"
	"github.com/src/main/app/log"
	"github.com/src/main/app/metrics"
)

type Pusher interface {
	SendMessage(message *queue.MessageDTO) error
}

type HTTPPusher struct {
	httpClient client.AppClient
}

type MessageDTO struct {
	ID        string `json:"MessageId,omitempty"`
	Message   string `json:"Message,omitempty"`
	Timestamp string `json:"Timestamp,omitempty"`
}

func NewHTTPPusher(httpClient client.AppClient) *HTTPPusher {
	return &HTTPPusher{
		httpClient: httpClient,
	}
}

func (h HTTPPusher) SendMessage(message *queue.MessageDTO) error {
	var messageDTO MessageDTO
	err := json.Unmarshal([]byte(message.Body), &messageDTO)
	if err != nil {
		log.Error(err)
		return err
	}

	requestBody := new(client.RequestBody)
	requestBody.ID = messageDTO.ID
	requestBody.Msg = messageDTO.Message
	requestBody.Timestamp = messageDTO.Timestamp

	log.Warnf("[pushing]: message id: %s, msg: %s, timestamp: %s", requestBody.ID, requestBody.Msg, requestBody.Timestamp)

	err = h.httpClient.PostMessage(requestBody)

	if err != nil {
		log.Errorf("[nack]   : message id: %s, msg: %s, timestamp: %s",
			requestBody.ID,
			requestBody.Msg,
			requestBody.Timestamp)

		metrics.Collector.IncrementCounter(metrics.PusherError)
		return err
	}

	log.Infof("[ack]    : message id: %s, msg: %s, timestamp: %s", requestBody.ID, requestBody.Msg, requestBody.Timestamp)
	metrics.Collector.IncrementCounter(metrics.PusherSuccess)

	return nil
}
