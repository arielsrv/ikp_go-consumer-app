package consumer

import (
	"context"
	"github.com/src/main/app/pusher"
	"github.com/src/main/app/queue"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type Config struct {
	// QueueClient URL to receive messages from.
	QueueURL string
	// Maximum workers that will independently receive messages from a queue.
	Workers int
	// Maximum messages that will be picked up by a worker in one-go.
	MaxMsg int
}

type Consumer struct {
	messageClient queue.MessageClient
	pusher        pusher.Pusher
	config        Config
}

func NewConsumer(messageClient queue.MessageClient, pusher pusher.Pusher, config Config) Consumer {
	return Consumer{
		messageClient: messageClient,
		pusher:        pusher,
		config:        config,
	}
}

func (c Consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(c.config.Workers)

	for i := 1; i <= c.config.Workers; i++ {
		go c.worker(ctx, wg, i)
	}

	wg.Wait()
}

func (c Consumer) worker(ctx context.Context, wg *sync.WaitGroup, workerId int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d: stopped\n", workerId)
			return
		default:
		}

		messages, err := c.messageClient.Receive(ctx, c.config.QueueURL, int64(c.config.MaxMsg))
		if err != nil {
			// Critical error!
			log.Printf("worker %d: receive error: %s\n", workerId, err.Error())
			continue
		}

		if len(messages) != 0 {
			c.asyncPop(ctx, messages)
		}
	}
}

func (c Consumer) asyncPop(ctx context.Context, messages []*sqs.Message) {
	wg := &sync.WaitGroup{}
	wg.Add(len(messages))

	for _, message := range messages {
		go func(message *sqs.Message) {
			defer wg.Done()
			c.pop(ctx, message)
		}(message)
	}

	wg.Wait()
}

func (c Consumer) pop(ctx context.Context, message *sqs.Message) {
	err := c.pusher.SendMessage(message)
	if err != nil {
		log.Printf("pusher error: %s, msg: %s\n", err.Error(), *message.Body)
	} else {
		err = c.messageClient.Delete(ctx, c.config.QueueURL, *message.ReceiptHandle)
		if err != nil {
			log.Printf("delete error: %s, msg: %s\n", err.Error(), *message.Body)
		}
	}
}
