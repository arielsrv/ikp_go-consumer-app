package queue

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
