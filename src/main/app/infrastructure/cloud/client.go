package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageClient interface {
	// Send a message to queue and returns its message ID.
	Send(ctx context.Context, sendRequest *SendRequest) (string, error)
	// Receive Long polls given amount of messages from a queue.
	Receive(ctx context.Context, queueURL string, maxMsg int64) ([]*sqs.Message, error)
	// Delete Deletes a message from a queue.
	Delete(ctx context.Context, queueURL, receiptHandle string) error
}
