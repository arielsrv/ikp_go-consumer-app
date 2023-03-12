package consumer

import (
	"context"
	"sync"
)

type TaskResolverType string

const (
	Sync  TaskResolverType = "sync"
	Async TaskResolverType = "async"
)

type TaskResolver[T any] struct {
	handlers map[string]ElementHandler[T]
}

func (r *TaskResolver[T]) Resolve(taskResolverType TaskResolverType) ElementHandler[T] {
	return r.handlers[string(taskResolverType)]
}

type ElementHandler[T any] interface {
	Process(ctx context.Context, elements []T, f func(ctx context.Context, element *T))
}

type asyncHandler[T any] struct {
}

func (handler asyncHandler[T]) Process(
	ctx context.Context,
	elements []T,
	f func(ctx context.Context, element *T),
) {
	wg := &sync.WaitGroup{}
	wg.Add(len(elements))

	for i := range elements {
		go func(element T) {
			defer wg.Done()
			f(ctx, &element)
		}(elements[i])
	}

	wg.Wait()
}

type syncHandler[T any] struct {
}

func (handler syncHandler[T]) Process(
	ctx context.Context,
	elements []T,
	f func(ctx context.Context, element *T),
) {
	for i := range elements {
		f(ctx, &elements[i])
	}
}
