package queue

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockSQS struct {
	sqsiface.SQSAPI
	messages map[string][]*sqs.Message
}

func getMockSQSClient() sqsiface.SQSAPI {
	return &MockSQS{
		messages: map[string][]*sqs.Message{},
	}
}

func (m *MockSQS) SendMessage(in *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	m.messages[*in.QueueUrl] = append(m.messages[*in.QueueUrl], &sqs.Message{
		Body: in.MessageBody,
	})
	return &sqs.SendMessageOutput{}, nil
}
func (m *MockSQS) ReceiveMessageWithContext(_ aws.Context, in *sqs.ReceiveMessageInput, _ ...request.Option) (*sqs.ReceiveMessageOutput, error) {
	if len(m.messages[*in.QueueUrl]) == 0 {
		return &sqs.ReceiveMessageOutput{}, nil
	}
	response := m.messages[*in.QueueUrl][0:1]
	m.messages[*in.QueueUrl] = m.messages[*in.QueueUrl][1:]
	return &sqs.ReceiveMessageOutput{
		Messages: response,
	}, nil
}

func TestNewClient(t *testing.T) {
	q := getMockSQSClient()
	queueURL := "https://queue.amazonaws.com/80398EXAMPLE/MyQueue"
	output, err := q.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    &queueURL,
	})
	assert.NoError(t, err)
	assert.NotNil(t, output)

	queue := NewClient(time.Second*5, q)
	message, err := queue.Receive(context.Background(), queueURL, 10)

	assert.NoError(t, err)
	assert.NotNil(t, message)
	assert.Len(t, message, 1)
	assert.NotNil(t, message[0])
	assert.NotNil(t, message[0].Body)
	assert.Equal(t, *message[0].Body, "Hello, world!")
}
