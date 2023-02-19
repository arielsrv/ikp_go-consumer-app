package pusher

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/helpers/types"
	"github.com/src/main/app/rest"
	"github.com/src/main/app/server/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) PostMessage(*rest.RequestBody) error {
	args := m.Called()
	return args.Error(0)
}

func TestHttpPusher_SendMessage(t *testing.T) {
	httpClient := new(MockHttpClient)
	httpPusher := NewHttpPusher(httpClient)

	httpClient.On("PostMessage").Return(nil)

	message := new(sqs.Message)
	message.Body = types.String("{\"MessageId\":\"123\", \"Message\": \"Hello world\"}")

	err := httpPusher.SendMessage(message)
	assert.NoError(t, err)
}

func TestHttpPusher_SendMessageErr(t *testing.T) {
	httpClient := new(MockHttpClient)
	httpPusher := NewHttpPusher(httpClient)

	httpClient.On("PostMessage").Return(errors.NewError(504, "gateway timeout"))

	message := new(sqs.Message)
	message.Body = types.String("{\"MessageId\":\"123\", \"Message\": \"Hello world\"}")

	err := httpPusher.SendMessage(message)
	assert.Error(t, err)
}

func TestHttpPusher_SendMessageParsingErr(t *testing.T) {
	httpClient := new(MockHttpClient)
	httpPusher := NewHttpPusher(httpClient)

	httpClient.On("PostMessage").Return(errors.NewError(504, "gateway timeout"))

	message := new(sqs.Message)
	message.Body = types.String("invalid message")

	err := httpPusher.SendMessage(message)
	assert.Error(t, err)
}
