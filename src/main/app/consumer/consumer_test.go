package consumer

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
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

	queue := queue.NewTestClient(time.Second * 5)
	output, err := queue.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, world!"),
		QueueUrl:    aws.String("https://queues.com/my-queue"),
	})
	assert.NoError(t, err)
	assert.NotNil(t, output)

	NewConsumer(queue, httpPusher, Config{
		QueueURL: "https://queues.com/my-queue",
		Workers:  1,
		MaxMsg:   1,
	}).Start(ctx)

	receiveMessageOutput, err := queue.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String("https://queues.com/my-queue"),
		MaxNumberOfMessages:   aws.Int64(int64(1)),
		WaitTimeSeconds:       aws.Int64(10),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})

	assert.NoError(t, err)
	assert.NotNil(t, receiveMessageOutput)
	assert.Nil(t, receiveMessageOutput.Messages)
}
