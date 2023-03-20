package resolvers

import (
	"context"
	"sync"
)

type AsyncResolver[T comparable] struct {
}

func (resolver AsyncResolver[T]) Process(ctx context.Context, elements []T, f func(ctx context.Context, element *T)) {
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
