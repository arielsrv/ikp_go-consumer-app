package queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/src/main/app/helpers/arrays"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
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
	MaxMsg   int
}

type MessageDTO struct {
	Body          string
	ReceiptHandle string
}

type Config struct {
	Name     string
	URL      string
	Parallel int
	Timeout  int
}

func NewClient(config Config, session *session.Session) (*Client, error) {
	if config.Parallel < 1 || config.Parallel > 10 {
		log.Errorf("receive argument: parallel valid values: 1 to 10: given %d", config.Parallel)
		return nil, errors.New("invalidad parallel value")
	}

	return &Client{
		timeout:  time.Millisecond * time.Duration(config.Timeout),
		SQSAPI:   sqs.New(session),
		QueueURL: config.URL,
		MaxMsg:   config.Parallel,
	}, nil
}

func (s Client) Receive(ctx context.Context) ([]MessageDTO, error) {
	var ctxTimeout = s.timeout + s.timeout/2

	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	receiveMessageOutput, err := s.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(s.QueueURL),
		MaxNumberOfMessages:   aws.Int64(int64(s.MaxMsg)),
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
