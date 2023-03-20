package consumer

import (
	"context"
	"sync"

	"github.com/ugurcsen/gods-generic/maps/hashmap"

	"github.com/src/main/app/consumer/resolvers"
	"github.com/src/main/app/queue"
)

type TaskResolverType string

const (
	Sync  TaskResolverType = "sync"
	Async TaskResolverType = "async"
)

type TaskResolver[T comparable] struct {
	handlers *hashmap.Map[TaskResolverType, ElementHandler[T]]
}

func (r *TaskResolver[T]) Resolve(taskResolverType TaskResolverType) ElementHandler[T] {
	value, found := r.handlers.Get(taskResolverType)
	if !found {
		return nil
	}
	return value
}

type ElementHandler[T comparable] interface {
	Process(ctx context.Context, elements []T, f func(ctx context.Context, element *T))
}

var (
	instance     sync.Once
	taskResolver = &TaskResolver[queue.MessageDTO]{}
)

func ProvideTaskResolver() *TaskResolver[queue.MessageDTO] {
	instance.Do(func() {
		taskResolver = &TaskResolver[queue.MessageDTO]{}
		taskResolver.handlers = hashmap.New[TaskResolverType, ElementHandler[queue.MessageDTO]]()
		taskResolver.handlers.Put(Sync, &resolvers.SyncResolver[queue.MessageDTO]{})
		taskResolver.handlers.Put(Async, &resolvers.AsyncResolver[queue.MessageDTO]{})
	})
	return taskResolver
}
