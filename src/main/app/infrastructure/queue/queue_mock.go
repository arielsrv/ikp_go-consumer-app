package queue

import (
	"container/list"
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ugurcsen/gods-generic/maps/hashmap"
)

type MockClient struct {
	AWSQueueService
	queues *hashmap.Map[string, *list.List]
}

type MockConfig struct {
	QueueURL string
	MaxMsg   int
	Queues   *hashmap.Map[string, *list.List]
}

func NewMockClient(config MockConfig) AWSQueueService {
	return AWSQueueService{
		AWSClient: &MockClient{
			queues: config.Queues,
		},
		QueueURL: config.QueueURL,
		MaxMsg:   config.MaxMsg,
	}
}

func (m *MockClient) ReceiveMessage(_ context.Context, params *sqs.ReceiveMessageInput, _ ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	queue, err := m.getQueue(aws.ToString(params.QueueUrl))
	if err != nil {
		return nil, err
	}

	var messages []types.Message
	if queue.Len() > 0 {
		i := 0
		for e := queue.Front(); e != nil && i < int(params.MaxNumberOfMessages); e = e.Next() {
			messages = append(messages, e.Value.(types.Message))
			i++
		}
		return &sqs.ReceiveMessageOutput{
			Messages: messages,
		}, nil
	}

	return &sqs.ReceiveMessageOutput{
		Messages: messages,
	}, nil
}

func (m *MockClient) DeleteMessage(_ context.Context, params *sqs.DeleteMessageInput, _ ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	queue, err := m.getQueue(aws.ToString(params.QueueUrl))
	if err != nil {
		return nil, err
	}

	if queue.Len() > 0 {
		receiptHandle := aws.ToString(params.ReceiptHandle)
		for e := queue.Front(); e != nil; e = e.Next() {
			if e.Value != nil {
				message, converted := e.Value.(types.Message)
				if !converted {
					return nil, fmt.Errorf("conversion error: %s", receiptHandle)
				}
				if aws.ToString(message.ReceiptHandle) == receiptHandle {
					queue.Remove(e)
					return &sqs.DeleteMessageOutput{}, nil
				}
			}
		}
		return nil, fmt.Errorf("delete error: not found %s", receiptHandle)
	}

	return nil, fmt.Errorf("delete error:  %s", aws.ToString(params.ReceiptHandle))
}

func (m *MockClient) GetQueueAttributes(_ context.Context, params *sqs.GetQueueAttributesInput, _ ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error) {
	queue, err := m.getQueue(aws.ToString(params.QueueUrl))
	if err != nil {
		return nil, err
	}

	return &sqs.GetQueueAttributesOutput{
		Attributes: map[string]string{
			"ApproximateNumberOfMessagesNotVisible": strconv.Itoa(queue.Len()),
		},
	}, nil
}

func (m *MockClient) getQueue(queueURL string) (*list.List, error) {
	queue, found := m.queues.Get(queueURL)
	if !found {
		return nil, fmt.Errorf("non-existing queue: %s", queueURL)
	}

	return queue, nil
}
