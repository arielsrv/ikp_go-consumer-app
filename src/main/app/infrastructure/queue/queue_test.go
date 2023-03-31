package queue_test

import (
	"container/list"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/src/main/app/config"
	"github.com/src/main/app/container"
	"github.com/src/main/app/infrastructure/queue"
	"github.com/stretchr/testify/assert"
	"github.com/ugurcsen/gods-generic/maps/hashmap"
)

func TestNewClient(t *testing.T) {
	queueClient, err := queue.NewClient(queue.Config{
		Name:     config.String("queue"),
		URL:      config.String("url"),
		Parallel: config.TryInt("queues.orders.parallel", 10),
		Timeout:  config.TryInt("queues.orders.timeout", 1000),
	}, container.ProvideAWSConfig())

	assert.NoError(t, err)
	assert.NotNil(t, queueClient)
}

func TestNewClientErr(t *testing.T) {
	queueClient, err := queue.NewClient(queue.Config{
		Name:     config.String("queues.orders.name"),
		URL:      config.String("url"),
		Parallel: config.TryInt("", 20),
		Timeout:  config.TryInt("queues.orders.timeout", 1000),
	}, container.ProvideAWSConfig())

	assert.Error(t, err)
	assert.Nil(t, queueClient)
}

func TestNewFakeClient(t *testing.T) {
	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack(types.Message{
		Body: aws.String("msg1"),
	})
	l.PushBack(types.Message{
		Body: aws.String("msg2"),
	})
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	actual, err := queueClient.Receive(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Len(t, actual, 2)
	assert.NotNil(t, actual[0])
	assert.NotNil(t, actual[0].Body)
	assert.Equal(t, "msg1", actual[0].String())
	assert.NotNil(t, actual[1])
	assert.NotNil(t, actual[1].Body)
	assert.Equal(t, "msg2", actual[1].String())
}

func TestNewFakeClientInvalidQueue(t *testing.T) {
	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack(types.Message{
		Body: aws.String("msg"),
	})
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: "invalid queue",
		MaxMsg:   2,
		Queues:   queues,
	})

	actual, err := queueClient.Receive(context.Background())
	assert.Error(t, err)
	assert.Nil(t, actual)

	actualCount, countErr := queueClient.Count(context.Background())
	assert.Error(t, countErr)
	assert.Nil(t, actualCount)
}

func TestNewFakeClientEmpty(t *testing.T) {
	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	actual, err := queueClient.Receive(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, actual)
}

func TestNewFakeClientCount(t *testing.T) {
	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack(types.Message{
		Body: aws.String("msg1"),
	})
	l.PushBack(types.Message{
		Body: aws.String("msg2"),
	})
	l.PushBack(types.Message{
		Body: aws.String("msg3"),
	})
	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	count, err := queueClient.Count(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, count)
	assert.Equal(t, 3, aws.ToInt(count))
}

func TestNewClientDelete(t *testing.T) {
	queueURL := "https://queues.com/my-queue"
	l := new(list.List)

	l.PushBack(types.Message{
		Body:          aws.String("msg1"),
		ReceiptHandle: aws.String("rpt1"),
	})
	l.PushBack(types.Message{
		Body:          aws.String("msg2"),
		ReceiptHandle: aws.String("rpt2"),
	})

	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	err := queueClient.Delete(context.Background(), "rpt1")
	assert.NoError(t, err)

	err = queueClient.Delete(context.Background(), "rpt1")
	assert.Error(t, err)
}

func TestNewClientDeletePtrErr(t *testing.T) {
	queueURL := "https://queues.com/my-queue"
	l := new(list.List)
	l.PushBack("invalid message")

	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	err := queueClient.Delete(context.Background(), "wait ...")
	assert.Error(t, err)
	assert.Errorf(t, err, "invalid conversion")
}

func TestNewClientDeleteEmptyErr(t *testing.T) {
	queueURL := "https://queues.com/my-queue"
	l := new(list.List)

	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: queueURL,
		MaxMsg:   2,
		Queues:   queues,
	})

	err := queueClient.Delete(context.Background(), "wait ...")
	assert.Error(t, err)
}

func TestNewClientDeleteNotFoundErr(t *testing.T) {
	queueURL := "https://queues.com/my-queue"
	l := new(list.List)

	queues := hashmap.New[string, *list.List]()
	queues.Put(queueURL, l)

	queueClient := queue.NewMockClient(queue.MockConfig{
		QueueURL: "invalid queue",
		MaxMsg:   2,
		Queues:   queues,
	})

	err := queueClient.Delete(context.Background(), "wait ...")
	assert.Error(t, err)
}
