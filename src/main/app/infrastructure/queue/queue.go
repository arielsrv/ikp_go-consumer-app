package queue

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/src/main/app/log"
)

type AWSClient interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
	GetQueueAttributes(ctx context.Context, params *sqs.GetQueueAttributesInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error)
}

type AWSQueueService struct {
	Timeout  time.Duration
	QueueURL string
	MaxMsg   int
	AWSClient
}

type Config struct {
	Name     string
	URL      string
	Parallel int
	Timeout  int
}

func NewClient(config Config, awsConfig aws.Config) (*AWSQueueService, error) {
	if config.Parallel < 1 || config.Parallel > 10 {
		log.Errorf("receive argument: parallel valid values: 1 to 10: given %d", config.Parallel)
		return nil, errors.New("invalidad parallel value")
	}

	return &AWSQueueService{
		Timeout:   time.Millisecond * time.Duration(config.Timeout),
		AWSClient: sqs.NewFromConfig(awsConfig),
		QueueURL:  config.URL,
		MaxMsg:    config.Parallel,
	}, nil
}

func (s AWSQueueService) Receive(ctx context.Context) ([]MessageDTO, error) {
	var ctxTimeout = s.Timeout + s.Timeout/2
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	receiveMessageOutput, err := s.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.QueueURL),
		MaxNumberOfMessages:   int32(s.MaxMsg),
		WaitTimeSeconds:       int32(s.Timeout.Seconds()),
		MessageAttributeNames: []string{"All"},
	})

	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	if receiveMessageOutput.Messages == nil || len(receiveMessageOutput.Messages) == 0 {
		return nil, nil
	}

	messages := make([]MessageDTO, len(receiveMessageOutput.Messages))
	for i, message := range receiveMessageOutput.Messages {
		messageDTO := new(MessageDTO)
		messageDTO.Body = aws.ToString(message.Body)
		messageDTO.ReceiptHandle = aws.ToString(message.ReceiptHandle)
		messages[i] = *messageDTO
	}

	return messages, nil
}

func (s AWSQueueService) Delete(ctx context.Context, receiptHandle string) error {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	if _, err := s.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.QueueURL),
		ReceiptHandle: aws.String(receiptHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (s AWSQueueService) Count(ctx context.Context) (*int, error) {
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()

	output, err := s.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(s.QueueURL),
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameAll},
	})

	if err != nil {
		return nil, fmt.Errorf("queue retrieving attribute error: %w", err)
	}

	value := output.Attributes["ApproximateNumberOfMessagesNotVisible"]
	count, err := strconv.Atoi(value)

	if err != nil {
		return nil, fmt.Errorf("parse attribute error: %w", err)
	}

	return aws.Int(count), nil
}
