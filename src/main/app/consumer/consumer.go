package consumer

import (
	"context"
	"sync"

	"github.com/src/main/app/log"

	"github.com/src/main/app/queue"

	"github.com/src/main/app/pusher"
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

var (
	instance     sync.Once
	taskResolver = &TaskResolver[queue.MessageDTO]{
		handlers: make(map[string]ElementHandler[queue.MessageDTO]),
	}
)

func ProvideTaskResolver() *TaskResolver[queue.MessageDTO] {
	instance.Do(func() {
		taskResolver = &TaskResolver[queue.MessageDTO]{
			handlers: make(map[string]ElementHandler[queue.MessageDTO]),
		}
		taskResolver.handlers[string(Sync)] = &syncHandler[queue.MessageDTO]{}
		taskResolver.handlers[string(Async)] = &asyncHandler[queue.MessageDTO]{}
	})
	return taskResolver
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
			// Critical error
			log.Errorf("worker %d: receive error: %s\n", workerID, err.Error())
			continue
		}

		if !isEmpty(messages) {
			c.taskResolver.
				Resolve(c.taskResolverType).
				Process(ctx, messages, c.sendAndDelete)
		}
	}
}

func isEmpty(messages []queue.MessageDTO) bool {
	return len(messages) == 0
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
