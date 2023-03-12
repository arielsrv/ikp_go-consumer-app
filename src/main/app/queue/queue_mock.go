package queue

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/src/main/app/helpers/types"
)

type MockSQS struct {
	sqsiface.SQSAPI
	messages map[string][]*sqs.Message
}

func NewTestClient(queueURL string) Client {
	return Client{
		timeout: time.Second * 5,
		SQSAPI: &MockSQS{
			messages: map[string][]*sqs.Message{},
		},
		QueueURL: queueURL,
	}
}

func (m *MockSQS) SendMessage(in *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	m.messages[*in.QueueUrl] = append(m.messages[*in.QueueUrl], &sqs.Message{
		Body:          in.MessageBody,
		ReceiptHandle: aws.String("receipt-handle"),
	})
	return &sqs.SendMessageOutput{}, nil
}

func (m *MockSQS) ReceiveMessageWithContext(
	_ aws.Context,
	in *sqs.ReceiveMessageInput,
	_ ...request.Option,
) (*sqs.ReceiveMessageOutput, error) {
	if len(m.messages[*in.QueueUrl]) == 0 {
		return &sqs.ReceiveMessageOutput{}, nil
	}
	response := m.messages[*in.QueueUrl][0:1]
	return &sqs.ReceiveMessageOutput{
		Messages: response,
	}, nil
}

func (m *MockSQS) DeleteMessageWithContext(
	_ aws.Context,
	in *sqs.DeleteMessageInput,
	_ ...request.Option,
) (*sqs.DeleteMessageOutput, error) {
	if len(m.messages[*in.QueueUrl]) == 0 {
		return nil, errors.New("empty queue")
	}

	for i := 0; i < len(m.messages[*in.QueueUrl]); i++ {
		if types.StringValue(m.messages[*in.QueueUrl][i].ReceiptHandle) == types.StringValue(in.ReceiptHandle) {
			m.messages[*in.QueueUrl] = m.messages[*in.QueueUrl][1:]
			return &sqs.DeleteMessageOutput{}, nil
		}
	}

	return nil, errors.New("delete error")
}
