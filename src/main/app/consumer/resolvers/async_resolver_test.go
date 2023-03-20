package resolvers_test

import (
	"context"
	"testing"
	"time"

	"github.com/src/main/app/consumer/resolvers"
	"github.com/stretchr/testify/assert"
	"github.com/ugurcsen/gods-generic/lists/arraylist"
)

func TestAsyncResolver_Process(t *testing.T) {
	resolver := resolvers.AsyncResolver[string]{}
	elements := arraylist.New[string]()
	elements.Add("1")
	elements.Add("2")
	elements.Add("3")

	i := 0
	startTime := time.Now()
	resolver.
		Process(context.Background(), elements.Values(), func(ctx context.Context, element *string) {
			i++
			time.Sleep(time.Duration(100) * time.Millisecond)
		})
	elapsedTime := time.Since(startTime)

	assert.Equal(t, 3, i)
	assert.GreaterOrEqual(t, elapsedTime, time.Duration(100)*time.Millisecond)
	assert.LessOrEqual(t, elapsedTime.Milliseconds(), time.Duration(110)*time.Millisecond)
}
