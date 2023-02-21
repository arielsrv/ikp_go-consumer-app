package queue

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	properties "github.com/src/main/app/config"
	"github.com/src/main/app/helpers/types"
	"github.com/src/main/app/log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageClient interface {
	Send(ctx context.Context, sendRequest *SendRequest) (string, error)
	Receive(ctx context.Context) ([]*sqs.Message, error)
	Delete(ctx context.Context, receiptHandle string) error
}

type SendRequest struct {
	QueueURL   string
	Body       string
	Attributes []Attribute
}

type Attribute struct {
	Key   string
	Value string
	Type  string
}

type Client struct {
	timeout time.Duration
	sqsiface.SQSAPI
	QueueUrl string
	MaxMsg   int64
}

type MockClient struct {
	Client
	messages map[string][]*sqs.Message
}

func NewClient(queueName string, parallel int, timeout int) (*Client, error) {
	session, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Credentials: credentials.
					NewStaticCredentials(
						properties.String("aws.id"),
						properties.String("aws.secret"), ""),
				Region:           aws.String(properties.String("aws.region")),
				Endpoint:         aws.String(properties.String("aws.url")),
				S3ForcePathStyle: aws.Bool(true),
			},
			Profile: properties.String("aws.profile"),
		},
	)

	if err != nil {
		log.Errorf("aws session error: %s", err)
		return nil, err
	}

	if parallel < 1 || parallel > 10 {
		log.Errorf("receive argument: queues.users.parallel valid values: 1 to 10: given %d", parallel)
		return nil, err
	}

	sqsClient := sqs.New(session)
	responseQueueUrl, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: types.String(queueName),
	})
	if err != nil {
		log.Errorf("sqs session error: %s", err)
		return nil, err
	}

	return &Client{
		timeout:  time.Millisecond * time.Duration(timeout),
		SQSAPI:   sqsClient,
		QueueUrl: types.StringValue(responseQueueUrl.QueueUrl),
		MaxMsg:   int64(parallel),
	}, nil
}

type MockSQS struct {
	sqsiface.SQSAPI
	messages map[string][]*sqs.Message
}

func NewTestClient(timeout time.Duration, queueUrl string) Client {
	return Client{
		timeout: timeout,
		SQSAPI: &MockSQS{
			messages: map[string][]*sqs.Message{},
		},
		QueueUrl: queueUrl,
	}
}

func (m *MockSQS) SendMessage(in *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	m.messages[*in.QueueUrl] = append(m.messages[*in.QueueUrl], &sqs.Message{
		Body:          in.MessageBody,
		ReceiptHandle: aws.String("receipt-handle"),
	})
	return &sqs.SendMessageOutput{}, nil
}

func (m *MockSQS) ReceiveMessageWithContext(_ aws.Context, in *sqs.ReceiveMessageInput, _ ...request.Option) (*sqs.ReceiveMessageOutput, error) {
	if len(m.messages[*in.QueueUrl]) == 0 {
		return &sqs.ReceiveMessageOutput{}, nil
	}
	response := m.messages[*in.QueueUrl][0:1]
	m.messages[*in.QueueUrl] = m.messages[*in.QueueUrl][1:]
	return &sqs.ReceiveMessageOutput{
		Messages: response,
	}, nil
}

func (m *MockSQS) DeleteMessageWithContext(_ aws.Context, in *sqs.DeleteMessageInput, _ ...request.Option) (*sqs.DeleteMessageOutput, error) {
	if len(m.messages[*in.QueueUrl]) == 0 {
		return &sqs.DeleteMessageOutput{}, nil
	}
	return &sqs.DeleteMessageOutput{}, nil
}

func (s Client) Send(ctx context.Context, sendRequest *SendRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	attributes := make(map[string]*sqs.MessageAttributeValue, len(sendRequest.Attributes))
	for _, attribute := range sendRequest.Attributes {
		attributes[attribute.Key] = &sqs.MessageAttributeValue{
			StringValue: aws.String(attribute.Value),
			DataType:    aws.String(attribute.Type),
		}
	}

	sendMessageOutput, err := s.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageAttributes: attributes,
		MessageBody:       aws.String(sendRequest.Body),
		QueueUrl:          aws.String(sendRequest.QueueURL),
	})
	if err != nil {
		return "", fmt.Errorf("send: %w", err)
	}

	return *sendMessageOutput.MessageId, nil
}

func (s Client) Receive(ctx context.Context) ([]*sqs.Message, error) {
	var waitTimeSeconds int64 = 10

	// Must always be above `WaitTimeSeconds` otherwise `ReceiveMessageWithContext`
	// trigger context timeout error.
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(waitTimeSeconds+5))
	defer cancel()

	receiveMessageOutput, err := s.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.QueueUrl),
		MaxNumberOfMessages:   aws.Int64(s.MaxMsg),
		WaitTimeSeconds:       aws.Int64(waitTimeSeconds),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	return receiveMessageOutput.Messages, nil
}

func (s Client) Delete(ctx context.Context, receiptHandle string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.QueueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
