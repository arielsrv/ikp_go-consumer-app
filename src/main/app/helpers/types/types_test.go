package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	actual := String("Hello world")
	assert.NotNil(t, actual)
}

func TestStringValue(t *testing.T) {
	value1 := String("Hello world")
	value2 := String("Hello world")

	actual1 := StringValue(value1)
	actual2 := StringValue(value2)
	actual3 := StringValue(nil)

	assert.True(t, actual1 == actual2)
	assert.Equal(t, "", actual3)
}
