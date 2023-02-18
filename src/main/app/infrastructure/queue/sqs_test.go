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
func (m *MockSQS) ReceiveMessageWithContext(ctx aws.Context, in *sqs.ReceiveMessageInput, _ ...request.Option) (*sqs.ReceiveMessageOutput, error) {
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
	q.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, World!"),
		QueueUrl:    &queueURL,
	})
	message, _ := q.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl: &queueURL,
	})
	assert.Equal(t, *message.Messages[0].Body, "Hello, World!")

	queue := NewClient(time.Second*5, getMockSQSClient())
	sendRequest := new(SendRequest)
	sendRequest.Body = "hello world"
	queue.Receive(context.Background(), queueURL, 10)
}
