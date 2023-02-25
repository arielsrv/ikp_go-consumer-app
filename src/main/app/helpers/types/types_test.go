package types_test

import (
	"testing"

	"github.com/src/main/app/helpers/types"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	actual := types.String("Hello world")
	assert.NotNil(t, actual)
}

func TestStringValue(t *testing.T) {
	value1 := types.String("Hello world")
	value2 := types.String("Hello world")

	actual1 := types.StringValue(value1)
	actual2 := types.StringValue(value2)
	actual3 := types.StringValue(nil)

	assert.True(t, actual1 == actual2)
	assert.Equal(t, "", actual3)
}
