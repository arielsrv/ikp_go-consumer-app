package pusher_test

import (
	"testing"

	"github.com/src/main/app/client"
	"github.com/src/main/app/infrastructure/queue"
	"github.com/src/main/app/pusher"
	"github.com/src/main/app/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) PostMessage(*client.RequestBody) error {
	args := m.Called()
	return args.Error(0)
}

func TestHttpPusher_SendMessage(t *testing.T) {
	httpClient := new(MockHTTPClient)
	httpPusher := pusher.NewHTTPPusher(httpClient)

	httpClient.On("PostMessage").Return(nil)

	message := new(queue.MessageDTO)
	message.Body = "{\"MessageId\":\"123\", \"Message\": \"Hello world\"}"

	err := httpPusher.SendMessage(message)
	assert.NoError(t, err)
}

func TestHttpPusher_SendMessageErr(t *testing.T) {
	httpClient := new(MockHTTPClient)
	httpPusher := pusher.NewHTTPPusher(httpClient)

	httpClient.On("PostMessage").Return(server.NewError(504, "gateway timeout"))

	message := new(queue.MessageDTO)
	message.Body = "{\"MessageId\":\"123\", \"Message\": \"Hello world\"}"

	err := httpPusher.SendMessage(message)
	assert.Error(t, err)
}

func TestHttpPusher_SendMessageParsingErr(t *testing.T) {
	httpClient := new(MockHTTPClient)
	httpPusher := pusher.NewHTTPPusher(httpClient)

	httpClient.On("PostMessage").Return(server.NewError(504, "gateway timeout"))

	message := new(queue.MessageDTO)
	message.Body = "invalid message"

	err := httpPusher.SendMessage(message)
	assert.Error(t, err)
}
