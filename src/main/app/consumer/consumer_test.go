package consumer_test

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/src/main/app/consumer"

	"github.com/src/main/app/queue"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPusher struct {
	mock.Mock
}

func (m *MockPusher) SendMessage(*sqs.Message) error {
	args := m.Called()
	return args.Error(0)
}

func TestNewConsumer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(nil)

	queueClient := queue.NewTestClient(time.Second*5, "https://queues.com/my-queue")
	output, err := queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, output)

	consumer.NewConsumer(
		consumer.Config{
			MessageClient: queueClient,
			Pusher:        httpPusher,
			Workers:       1}).
		Start(ctx)

	receiveMessageOutput, err := queueClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   aws.Int64(int64(1)),
		WaitTimeSeconds:       aws.Int64(10),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.Nil(t, receiveMessageOutput.Messages)
}

func TestNewConsumerErr(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(errors.New("http client error"))

	queueClient := queue.NewTestClient(time.Second*5, "https://queues.com/my-queue")
	output, err := queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, output)

	consumer.NewConsumer(
		consumer.Config{
			MessageClient: queueClient,
			Pusher:        httpPusher,
			Workers:       1}).
		Start(ctx)

	receiveMessageOutput, err := queueClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   aws.Int64(int64(1)),
		WaitTimeSeconds:       aws.Int64(10),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.Nil(t, receiveMessageOutput.Messages)
	log.Println()
}
