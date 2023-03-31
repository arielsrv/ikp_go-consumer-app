package consumer_test

import (
	"container/list"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/src/main/app/consumer"
	"github.com/src/main/app/container"
	"github.com/src/main/app/infrastructure/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ugurcsen/gods-generic/maps/hashmap"
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

	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack(types.Message{
		Body: aws.String("msg"),
	})
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	receiveMessageOutput, err := queueClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   int32(1),
		WaitTimeSeconds:       int32(10),
		MessageAttributeNames: []string{"All"},
	})

	consumer.NewConsumer(
		consumer.Config{
			QueueService:     queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: consumer.Sync,
		}, container.ProvideConsumerService()).
		Start(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.NotNil(t, receiveMessageOutput.Messages)
	assert.NotNil(t, receiveMessageOutput.Messages[0])
	assert.Equal(t, "msg", aws.ToString(receiveMessageOutput.Messages[0].Body))
}

func TestNewConsumerAsyncEmpty(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(nil)

	queueURL := "https://queues.com/my-queue"
	l := new(list.List)

	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	receiveMessageOutput, err := queueClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueURL),
		MaxNumberOfMessages:   int32(1),
		WaitTimeSeconds:       int32(10),
		MessageAttributeNames: []string{"All"},
	})

	consumer.NewConsumer(
		consumer.Config{
			QueueService:     queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: consumer.Sync,
		}, container.ProvideConsumerService()).
		Start(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.Nil(t, receiveMessageOutput.Messages)
}

func TestNewConsumerSync(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(nil)

	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack(types.Message{
		Body: aws.String("msg"),
	})
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	receiveMessageOutput, err := queueClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   int32(1),
		WaitTimeSeconds:       int32(10),
		MessageAttributeNames: []string{"All"},
	})

	consumer.NewConsumer(
		consumer.Config{
			QueueService:     queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: consumer.Async,
		}, container.ProvideConsumerService()).
		Start(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.NotNil(t, receiveMessageOutput.Messages)
	assert.NotNil(t, receiveMessageOutput.Messages[0])
	assert.Equal(t, "msg", aws.ToString(receiveMessageOutput.Messages[0].Body))
}

func TestNewConsumerSyncErr(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(errors.New("internal server error"))

	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack(types.Message{
		Body: aws.String("msg"),
	})
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	_, _ = queueClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       10,
		MessageAttributeNames: []string{"All"},
	})

	consumer.NewConsumer(
		consumer.Config{
			QueueService:     queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: consumer.Async,
		}, container.ProvideConsumerService()).
		Start(ctx)

	receiveMessageOutput, err := queueClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       10,
		MessageAttributeNames: []string{"All"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.NotNil(t, receiveMessageOutput.Messages)
	assert.NotNil(t, receiveMessageOutput.Messages[0])
	assert.Equal(t, "msg", aws.ToString(receiveMessageOutput.Messages[0].Body))
}

func TestNewConsumerSyncResolveErr(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(errors.New("internal server error"))

	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack(types.Message{
		Body: aws.String("msg"),
	})
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	_, _ = queueClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       10,
		MessageAttributeNames: []string{"All"},
	})

	consumer.NewConsumer(
		consumer.Config{
			QueueService:     queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: "invalid resolver",
		}, container.ProvideConsumerService()).
		Start(ctx)

	receiveMessageOutput, err := queueClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       10,
		MessageAttributeNames: []string{"All"},
	})

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.NotNil(t, receiveMessageOutput.Messages)
	assert.NotNil(t, receiveMessageOutput.Messages[0])
	assert.Equal(t, "msg", aws.ToString(receiveMessageOutput.Messages[0].Body))
}

func TestNewConsumerStopped(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(500))
	defer cancel()

	httpPusher := new(MockPusher)
	httpPusher.On("SendMessage").Return(nil)

	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack(types.Message{
		Body: aws.String("msg"),
	})
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	receiveMessageOutput, err := queueClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       10,
		MessageAttributeNames: []string{"All"},
	})

	consumerService := container.ProvideConsumerService()
	assert.NoError(t, consumerService.Stop())

	consumer.NewConsumer(
		consumer.Config{
			QueueService:     queueClient,
			Pusher:           httpPusher,
			Workers:          1,
			TaskResolverType: consumer.Sync,
		}, consumerService).
		Start(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.NotNil(t, receiveMessageOutput.Messages)
	assert.NotNil(t, receiveMessageOutput.Messages[0])
	assert.Equal(t, "msg", aws.ToString(receiveMessageOutput.Messages[0].Body))
}
