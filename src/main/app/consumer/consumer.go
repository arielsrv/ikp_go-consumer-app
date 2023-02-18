package consumer

import (
	"context"
	"github.com/src/main/app/infrastructure/cloud"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type Type string

const (
	// SyncConsumer Consumers to consume messages one by one.
	// A single goroutine handles all messages.
	// Progression is slower and requires less system resource.
	// Ideal for quiet/non-critical queues.
	SyncConsumer Type = "blocking"
	// AsyncConsumer Consumers to consume messages at the same time.
	// Runs an individual goroutine per message.
	// Progression is faster and requires more system resource.
	// Ideal for busy/critical queues.
	AsyncConsumer Type = "non-blocking"
)

type Config struct {
	// Instructs whether to consume messages come from a worker synchronously or asynchronous.
	Type Type
	// Queue URL to receive messages from.
	QueueURL string
	// Maximum workers that will independently receive messages from a queue.
	MaxWorker int
	// Maximum messages that will be picked up by a worker in one-go.
	MaxMsg int
}

type Consumer struct {
	client cloud.MessageClient
	config Config
}

func NewConsumer(client cloud.MessageClient, config Config) Consumer {
	return Consumer{
		client: client,
		config: config,
	}
}

func (c Consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(c.config.MaxWorker)

	for i := 1; i <= c.config.MaxWorker; i++ {
		go c.worker(ctx, wg, i)
	}

	wg.Wait()
}

func (c Consumer) worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	log.Printf("worker %d: started\n", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d: stopped\n", id)
			return
		default:
		}

		messages, err := c.client.Receive(ctx, c.config.QueueURL, int64(c.config.MaxMsg))
		if err != nil {
			// Critical error!
			log.Printf("worker %d: receive error: %s\n", id, err.Error())
			continue
		}

		if len(messages) == 0 {
			continue
		}

		if c.config.Type == SyncConsumer {
			c.sync(ctx, messages)
		} else {
			c.async(ctx, messages)
		}
	}
}

func (c Consumer) sync(ctx context.Context, messages []*sqs.Message) {
	for _, message := range messages {
		c.consume(ctx, message)
	}
}

func (c Consumer) async(ctx context.Context, messages []*sqs.Message) {
	wg := &sync.WaitGroup{}
	wg.Add(len(messages))

	for _, message := range messages {
		go func(message *sqs.Message) {
			defer wg.Done()
			c.consume(ctx, message)
		}(message)
	}

	wg.Wait()
}

func (c Consumer) consume(ctx context.Context, message *sqs.Message) {
	log.Println(*message.Body)
	if err := c.client.Delete(ctx, c.config.QueueURL, *message.ReceiptHandle); err != nil {
		// Critical error!
		log.Printf("delete error: %s\n", err.Error())
	}
}
