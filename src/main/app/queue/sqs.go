package queue

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	properties "github.com/src/main/app/config"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageClient interface {
	// Send a message to queue and returns its message ID.
	Send(ctx context.Context, sendRequest *SendRequest) (string, error)
	// Receive Long polls given amount of messages from a queue.
	Receive(ctx context.Context) ([]*sqs.Message, error)
	// Delete Deletes a message from a queue.
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

func NewClient(queueUrl string) Client {
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
		log.Fatalf("aws session error: %s", err)
	}

	maxMsg := properties.TryInt("consumers.users.workers.messages", 10)
	if maxMsg < 1 || maxMsg > 10 {
		fmt.Errorf("receive argument: msgMax valid values: 1 to 10: given %d", maxMsg)
	}

	return Client{
		timeout:  time.Millisecond * time.Duration(properties.TryInt("context.timeout", 1000)),
		SQSAPI:   sqs.New(session),
		QueueUrl: queueUrl,
		MaxMsg:   int64(maxMsg),
	}
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
