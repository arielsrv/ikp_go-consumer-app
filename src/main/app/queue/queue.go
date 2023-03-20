package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/src/main/app/helpers/arrays"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	properties "github.com/src/main/app/config"
	"github.com/src/main/app/helpers/types"
	"github.com/src/main/app/log"
)

type MessageClient interface {
	Receive(ctx context.Context) ([]MessageDTO, error)
	Delete(ctx context.Context, receiptHandle string) error
}

type Client struct {
	timeout time.Duration
	sqsiface.SQSAPI
	QueueURL string
	MaxMsg   int64
}

type MessageDTO struct {
	Body          string
	ReceiptHandle string
}

type Config struct {
	QueueName string
	Parallel  int
	Timeout   int
}

func NewClient(config Config) (*Client, error) {
	if config.Parallel < 1 || config.Parallel > 10 {
		log.Errorf("receive argument: parallel valid values: 1 to 10: given %d", config.Parallel)
		return nil, errors.New("invalidad parallel value")
	}

	awsSession, err := session.NewSessionWithOptions(
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

	sqsClient := sqs.New(awsSession)
	responseQueueURL, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: types.String(config.QueueName),
	})
	if err != nil {
		log.Errorf("sqs session error: %s", err)
		return nil, err
	}

	return &Client{
		timeout:  time.Millisecond * time.Duration(config.Timeout),
		SQSAPI:   sqsClient,
		QueueURL: types.StringValue(responseQueueURL.QueueUrl),
		MaxMsg:   int64(config.Parallel),
	}, nil
}

func (s Client) Receive(ctx context.Context) ([]MessageDTO, error) {
	var ctxTimeout = s.timeout + s.timeout/2

	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	receiveMessageOutput, err := s.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.QueueURL),
		MaxNumberOfMessages:   aws.Int64(s.MaxMsg),
		WaitTimeSeconds:       aws.Int64(int64(s.timeout.Seconds())),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	if !arrays.IsEmpty(receiveMessageOutput.Messages) {
		messages := make([]MessageDTO, len(receiveMessageOutput.Messages))
		for i, message := range receiveMessageOutput.Messages {
			messageDTO := new(MessageDTO)
			messageDTO.Body = types.StringValue(message.Body)
			messageDTO.ReceiptHandle = types.StringValue(message.ReceiptHandle)
			messages[i] = *messageDTO
		}

		return messages, nil
	}

	return nil, nil
}

func (s Client) Delete(ctx context.Context, receiptHandle string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.QueueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
