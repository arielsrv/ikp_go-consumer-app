package queue_test

import (
	"context"
	"testing"

	"github.com/src/main/app/container"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/config"
	"github.com/src/main/app/queue"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	queueClient, err := queue.NewClient(queue.Config{
		Name:     config.String("queue"),
		URL:      config.String("url"),
		Parallel: config.TryInt("queues.orders.parallel", 10),
		Timeout:  config.TryInt("queues.orders.timeout", 1000),
	}, container.ProvideAWSSession())

	assert.NoError(t, err)
	assert.NotNil(t, queueClient)
}

func TestNewClientErr(t *testing.T) {
	queueClient, err := queue.NewClient(queue.Config{
		Name:     config.String("queues.orders.name"),
		URL:      config.String("url"),
		Parallel: config.TryInt("", 20),
		Timeout:  config.TryInt("queues.orders.timeout", 1000),
	}, container.ProvideAWSSession())

	assert.Error(t, err)
	assert.Nil(t, queueClient)
}

func TestNewFakeClient(t *testing.T) {
	queueClient := queue.NewTestClient("https://queues.com/my-queue")

	output, err := queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})

	assert.NoError(t, err)
	assert.NotNil(t, output)

	actual, err := queueClient.Receive(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 1)
	assert.NotNil(t, actual[0])
	assert.NotNil(t, actual[0].Body)
	assert.Equal(t, actual[0].Body, "Hello, world!")
}

func TestNewFakeClientEmpty(t *testing.T) {
	queueClient := queue.NewTestClient("https://queues.com/my-queue")

	actual, err := queueClient.Receive(context.Background())

	assert.NoError(t, err)
	assert.Nil(t, actual)
}

func TestNewClientDelete(t *testing.T) {
	queueClient := queue.NewTestClient("https://queues.com/my-queue")

	output, err := queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})

	assert.NoError(t, err)
	assert.NotNil(t, output)

	actual, err := queueClient.Receive(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 1)
	assert.NotNil(t, actual[0])
	assert.NotNil(t, actual[0].Body)
	assert.Equal(t, actual[0].Body, "Hello, world!")

	err = queueClient.Delete(context.Background(), actual[0].ReceiptHandle)
	assert.NoError(t, err)

	err = queueClient.Delete(context.Background(), actual[0].ReceiptHandle)
	assert.Error(t, err)
}

func TestNewClientDelete_NotFound(t *testing.T) {
	queueClient := queue.NewTestClient("https://queues.com/my-queue")

	output, err := queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})

	assert.NoError(t, err)
	assert.NotNil(t, output)

	output, err = queueClient.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})

	assert.NoError(t, err)
	assert.NotNil(t, output)

	actual, err := queueClient.Receive(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 1)
	assert.NotNil(t, actual[0])
	assert.NotNil(t, actual[0].Body)
	assert.Equal(t, actual[0].Body, "Hello, world!")

	err = queueClient.Delete(context.Background(), actual[0].ReceiptHandle)
	assert.NoError(t, err)

	err = queueClient.Delete(context.Background(), "not found message")
	assert.Error(t, err)
}
