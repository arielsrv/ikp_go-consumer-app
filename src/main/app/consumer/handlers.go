package consumer

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type HandlerType string

const (
	Sync  HandlerType = "sync"
	Async HandlerType = "async"
)

var (
	instance        sync.Once
	handlerResolver = &HandlerResolver{
		handlers: make(map[string]MessageHandler),
	}
)

func ProvideHandlerResolver() *HandlerResolver {
	instance.Do(func() {
		handlerResolver = &HandlerResolver{
			handlers: make(map[string]MessageHandler),
		}
		handlerResolver.handlers[string(Sync)] = &syncHandler{}
		handlerResolver.handlers[string(Async)] = &asyncHandler{}
	})
	return handlerResolver
}

type HandlerResolver struct {
	handlers map[string]MessageHandler
}

func (r *HandlerResolver) resolve(handlerType HandlerType) MessageHandler {
	return r.handlers[string(handlerType)]
}

type MessageHandler interface {
	process(ctx context.Context, messages []*sqs.Message, handleMessage func(ctx context.Context, message *sqs.Message))
}

type asyncHandler struct {
}

func (handler asyncHandler) process(ctx context.Context, messages []*sqs.Message, handle func(ctx context.Context, message *sqs.Message)) {
	wg := &sync.WaitGroup{}
	wg.Add(len(messages))

	for _, message := range messages {
		go func(message *sqs.Message) {
			defer wg.Done()
			handle(ctx, message)
		}(message)
	}

	wg.Wait()
}

type syncHandler struct {
}

func (handler syncHandler) process(ctx context.Context, messages []*sqs.Message, handleMessage func(ctx context.Context, message *sqs.Message)) {
	for _, message := range messages {
		handleMessage(ctx, message)
	}
}
