package pusher

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/clients"
)

type Pusher interface {
	SendMessage(message *sqs.Message) error
}

type HttpPusher struct {
	httpClient clients.IHttpClient
}

type MessageDTO struct {
	Message string
}

func NewHttpPusher(httpClient clients.IHttpClient) *HttpPusher {
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
	requestBody := new(clients.RequestBody)
	requestBody.Msg = messageDTO.Message
	err = h.httpClient.PostMessage(requestBody)
	if err != nil {
		return err
	}
	return nil
}
