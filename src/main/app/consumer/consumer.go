package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/src/main/app/helpers/arrays"
	"github.com/src/main/app/infrastructure/queue"
	"github.com/src/main/app/log"
	"github.com/src/main/app/metrics"
	"github.com/src/main/app/model"
	"github.com/src/main/app/pusher"
	"github.com/src/main/app/services"
)

type Consumer struct {
	queueService     queue.Service
	pusher           pusher.Pusher
	workers          int
	taskResolverType TaskResolverType
	taskResolver     *TaskResolver[queue.MessageDTO]
	consumerService  services.IConsumerService
}

type Config struct {
	QueueService     queue.Service
	Pusher           pusher.Pusher
	Workers          int
	TaskResolverType TaskResolverType
}

func NewConsumer(config Config, consumerService services.IConsumerService) Consumer {
	return Consumer{
		queueService:     config.QueueService,
		pusher:           config.Pusher,
		workers:          config.Workers,
		taskResolverType: config.TaskResolverType,
		taskResolver:     ProvideTaskResolver(),
		consumerService:  consumerService,
	}
}

func (c Consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(c.workers + 1)

	for i := 0; i < c.workers; i++ {
		go c.worker(ctx, wg, i)
	}

	go c.collectMetrics(ctx, wg, c.workers)

	wg.Wait()
}

func (c Consumer) collectMetrics(ctx context.Context, wg *sync.WaitGroup, currentWorkers int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Infof("worker %d: stopped\n", currentWorkers)
			return
		default:
		}

		metrics.Collector.Record(metrics.CurrentWorkers, currentWorkers)

		approximateNumberOfMessages, err := c.queueService.Count(ctx)
		if err != nil {
			log.Warnf("metrics approximateNumberOfMessages error: %s", err.Error())
			time.Sleep(time.Millisecond * 1000)
			continue
		}

		metrics.Collector.Record(metrics.ApproximateNumberOfMessages, aws.ToInt(approximateNumberOfMessages))
		time.Sleep(time.Millisecond * 1000)
	}
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

		if c.consumerService.GetAppStatus().Status == model.Stopped {
			time.Sleep(time.Millisecond * 1000)
			continue
		}

		messages, err := c.queueService.Receive(ctx)
		if err != nil {
			log.Errorf("worker %d: critical receive error: %s\n", workerID, err.Error())
			time.Sleep(time.Millisecond * 5000)
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
		err = c.queueService.Delete(ctx, message.ReceiptHandle)
		if err != nil {
			log.Errorf("delete error: %s, msg: %s\n", err.Error(), message.Body)
		}
	}
}
