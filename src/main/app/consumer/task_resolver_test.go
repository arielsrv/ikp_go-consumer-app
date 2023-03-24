package consumer_test

import (
	"testing"

	"github.com/src/main/app/consumer"
	"github.com/stretchr/testify/assert"
)

func TestProvideTaskResolver(t *testing.T) {
	taskResolver := consumer.ProvideTaskResolver()
	actual, err := taskResolver.Resolve("mixed")

	assert.Error(t, err)
	assert.Equal(t, "invalid task resolver type: mixed", err.Error())
	assert.Nil(t, actual)
}
