package consumer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/consumer"
	"github.com/src/main/app/helpers/types"
	"github.com/src/main/app/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPusher struct {
	mock.Mock
}

func (m *MockPusher) SendMessage(*queue.MessageDTO) error {
	args := m.Called()
	return args.Error(0)
}

func TestNewConsumerAsync(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(nil)

	queueClient := queue.NewTestClient("https://queues.com/my-queue")
	output, err := queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, output)

	receiveMessageOutput, err := queueClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   aws.Int64(int64(1)),
		WaitTimeSeconds:       aws.Int64(10),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	consumer.NewConsumer(
		consumer.Config{
			MessageClient:    queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: consumer.Sync,
		}).
		Start(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.NotNil(t, receiveMessageOutput.Messages)
	assert.NotNil(t, receiveMessageOutput.Messages[0])
	assert.Equal(t, receiveMessageOutput.Messages[0].Body, types.String("Hello, world!"))
}

func TestNewConsumerSync(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(nil)

	queueClient := queue.NewTestClient("https://queues.com/my-queue")
	output, err := queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, output)

	receiveMessageOutput, err := queueClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   aws.Int64(int64(1)),
		WaitTimeSeconds:       aws.Int64(10),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	consumer.NewConsumer(
		consumer.Config{
			MessageClient:    queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: consumer.Async,
		}).
		Start(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.NotNil(t, receiveMessageOutput.Messages)
	assert.NotNil(t, receiveMessageOutput.Messages[0])
	assert.Equal(t, receiveMessageOutput.Messages[0].Body, types.String("Hello, world!"))
}

func TestNewConsumerSyncErr(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(errors.New("internal server error"))

	queueClient := queue.NewTestClient("https://queues.com/my-queue")
	output, err := queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, output)

	_, _ = queueClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   aws.Int64(int64(1)),
		WaitTimeSeconds:       aws.Int64(10),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	consumer.NewConsumer(
		consumer.Config{
			MessageClient:    queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: consumer.Async,
		}).
		Start(ctx)

	receiveMessageOutput, err := queueClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   aws.Int64(int64(1)),
		WaitTimeSeconds:       aws.Int64(10),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.NotNil(t, receiveMessageOutput.Messages)
	assert.NotNil(t, receiveMessageOutput.Messages[0])
	assert.Equal(t, receiveMessageOutput.Messages[0].Body, types.String("Hello, world!"))
}
