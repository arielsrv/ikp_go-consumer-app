package consumer

import (
	"context"
	"sync"

	"github.com/src/main/app/log"

	"github.com/src/main/app/queue"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/src/main/app/pusher"
)

type Consumer struct {
	messageClient queue.MessageClient
	pusher        pusher.Pusher
	workers       int
}

type Config struct {
	MessageClient queue.MessageClient
	Pusher        pusher.Pusher
	Workers       int
}

func NewConsumer(config Config) Consumer {
	return Consumer{
		messageClient: config.MessageClient,
		pusher:        config.Pusher,
		workers:       config.Workers,
	}
}

func (c Consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(c.workers)

	for i := 1; i <= c.workers; i++ {
		go c.worker(ctx, wg, i)
	}

	wg.Wait()
}

func (c Consumer) worker(ctx context.Context, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Infof("worker %d: stopped\n", workerID)
			return
		default:
		}

		messages, err := c.messageClient.Receive(ctx)
		if err != nil {
			// Critical error
			log.Errorf("worker %d: receive error: %s\n", workerID, err.Error())
			continue
		}

		if len(messages) != 0 {
			c.iterateAndPopAsync(ctx, messages)
		}
	}
}

func (c Consumer) iterateAndPopAsync(ctx context.Context, messages []*sqs.Message) {
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
		log.Errorf("pusher error: %s, msg: %s\n", err.Error(), *message.Body)
	} else {
		err = c.messageClient.Delete(ctx, *message.ReceiptHandle)
		if err != nil {
			log.Errorf("delete error: %s, msg: %s\n", err.Error(), *message.Body)
		}
	}
}
