package queue_test

import (
	"context"
	"testing"
	"time"

	"github.com/src/main/app/queue"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	queueClient := queue.NewTestClient(time.Second*5, "https://queues.com/my-queue")

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
	assert.Equal(t, *actual[0].Body, "Hello, world!")
}

func TestNewClientDelete(t *testing.T) {
	queueClient := queue.NewTestClient(time.Second*5, "https://queues.com/my-queue")

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
	assert.Equal(t, *actual[0].Body, "Hello, world!")

	err = queueClient.Delete(context.Background(), "receipt-handle")
	assert.NoError(t, err)
}
