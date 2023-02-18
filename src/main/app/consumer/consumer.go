package consumer

import (
	"context"
	"encoding/json"
	"github.com/src/main/app/clients"
	"github.com/src/main/app/infrastructure/cloud"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type Config struct {
	// Queue URL to receive messages from.
	QueueURL string
	// Maximum workers that will independently receive messages from a queue.
	Workers int
	// Maximum messages that will be picked up by a worker in one-go.
	MaxMsg int
}

type Consumer struct {
	messageClient cloud.MessageClient
	httpClient    clients.Client
	config        Config
}

func NewConsumer(messageClient cloud.MessageClient, httpClient clients.Client, config Config) Consumer {
	return Consumer{
		messageClient: messageClient,
		httpClient:    httpClient,
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
			c.iterateAsync(ctx, messages)
		}
	}
}

func (c Consumer) iterateAsync(ctx context.Context, messages []*sqs.Message) {
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

type AWSMessageDTO struct {
	Message string
}

func (c Consumer) pop(ctx context.Context, message *sqs.Message) {
	var awsMessage AWSMessageDTO
	err := json.Unmarshal([]byte(*message.Body), &awsMessage)
	if err != nil {
		log.Printf("invalid message error: %s, msg: %s\n", err.Error(), *message.Body)
	} else {
		requestBody := new(clients.RequestBody)
		requestBody.Msg = awsMessage.Message
		err = c.httpClient.PostMessage(requestBody)
		if err != nil {
			log.Printf("pusher error: %s, msg: %s\n", err.Error(), *message.Body)
		} else {
			err = c.messageClient.Delete(ctx, c.config.QueueURL, *message.ReceiptHandle)
			if err != nil {
				log.Printf("delete error: %s, msg: %s\n", err.Error(), *message.Body)
			}
		}
	}
}
