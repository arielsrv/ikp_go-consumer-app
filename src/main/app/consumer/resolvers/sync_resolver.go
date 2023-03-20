package resolvers

import "context"

type SyncResolver[T comparable] struct {
}

func (r SyncResolver[T]) Process(ctx context.Context, elements []T, f func(ctx context.Context, element *T)) {
	for i := range elements {
		f(ctx, &elements[i])
	}
}
