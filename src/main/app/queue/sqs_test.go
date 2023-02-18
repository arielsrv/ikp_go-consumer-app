package queue

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
	queue := NewTestClient(time.Second * 5)

	output, err := queue.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})

	assert.NoError(t, err)
	assert.NotNil(t, output)

	actual, err := queue.Receive(context.Background(), "https://queues.com/my-queue", 10)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 1)
	assert.NotNil(t, actual[0])
	assert.NotNil(t, actual[0].Body)
	assert.Equal(t, *actual[0].Body, "Hello, world!")
}
