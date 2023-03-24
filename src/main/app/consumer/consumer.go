package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/src/main/app/helpers/arrays"
	"github.com/src/main/app/log"
	"github.com/src/main/app/pusher"
	"github.com/src/main/app/queue"
)

type Consumer struct {
	messageClient    queue.MessageClient
	pusher           pusher.Pusher
	workers          int
	taskResolverType TaskResolverType
	taskResolver     *TaskResolver[queue.MessageDTO]
}

type Config struct {
	MessageClient    queue.MessageClient
	Pusher           pusher.Pusher
	Workers          int
	TaskResolverType TaskResolverType
}

func NewConsumer(config Config) Consumer {
	return Consumer{
		messageClient:    config.MessageClient,
		pusher:           config.Pusher,
		workers:          config.Workers,
		taskResolverType: config.TaskResolverType,
		taskResolver:     ProvideTaskResolver(),
	}
}

func (c Consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(c.workers)

	// workers.Each(func() { c.worker
	for i := 0; i < c.workers; i++ {
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
			log.Errorf("worker %d: critical receive error: %s\n", workerID, err.Error())
			time.Sleep(time.Millisecond * 1000)
			continue
		}

		if !arrays.IsEmpty(messages) {
			resolver, resolverErr := c.taskResolver.Resolve(c.taskResolverType)
			if resolverErr != nil {
				log.Errorf("worker %d: critical resolver error: %s\n", workerID, resolverErr.Error())
				time.Sleep(time.Millisecond * 1000)
				continue
			}
			resolver.Process(ctx, messages, c.sendAndDelete)
		}
	}
}

func (c Consumer) sendAndDelete(ctx context.Context, message *queue.MessageDTO) {
	err := c.pusher.SendMessage(message)
	if err != nil {
		log.Errorf("pusher error: %s, msg: %s\n", err.Error(), message.Body)
	} else {
		err = c.messageClient.Delete(ctx, message.ReceiptHandle)
		if err != nil {
			log.Errorf("delete error: %s, msg: %s\n", err.Error(), message.Body)
		}
	}
}
