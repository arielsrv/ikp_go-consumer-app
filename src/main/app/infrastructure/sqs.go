package infrastructure

import (
	"context"
	"fmt"
	"github.com/src/main/app/infrastructure/cloud"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var _ cloud.MessageClient = SQS{}

type SQS struct {
	timeout time.Duration
	client  *sqs.SQS
}

func NewSQS(session *session.Session, timeout time.Duration) SQS {
	return SQS{
		timeout: timeout,
		client:  sqs.New(session),
	}
}

func (s SQS) Send(ctx context.Context, req *cloud.SendRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	attributes := make(map[string]*sqs.MessageAttributeValue, len(req.Attributes))
	for _, attribute := range req.Attributes {
		attributes[attribute.Key] = &sqs.MessageAttributeValue{
			StringValue: aws.String(attribute.Value),
			DataType:    aws.String(attribute.Type),
		}
	}

	sendMessageOutput, err := s.client.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageAttributes: attributes,
		MessageBody:       aws.String(req.Body),
		QueueUrl:          aws.String(req.QueueURL),
	})
	if err != nil {
		return "", fmt.Errorf("send: %w", err)
	}

	return *sendMessageOutput.MessageId, nil
}

func (s SQS) Receive(ctx context.Context, queueURL string, maxMsg int64) ([]*sqs.Message, error) {
	if maxMsg < 1 || maxMsg > 10 {
		return nil, fmt.Errorf("receive argument: msgMax valid values: 1 to 10: given %d", maxMsg)
	}

	var waitTimeSeconds int64 = 10

	// Must always be above `WaitTimeSeconds` otherwise `ReceiveMessageWithContext`
	// trigger context timeout error.
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(waitTimeSeconds+5))
	defer cancel()

	receiveMessageOutput, err := s.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueURL),
		MaxNumberOfMessages:   aws.Int64(maxMsg),
		WaitTimeSeconds:       aws.Int64(waitTimeSeconds),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	return receiveMessageOutput.Messages, nil
}

func (s SQS) Delete(ctx context.Context, queueURL, rcvHandle string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(rcvHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
