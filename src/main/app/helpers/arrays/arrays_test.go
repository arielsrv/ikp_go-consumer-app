package arrays_test

import (
	"testing"

	"github.com/src/main/app/helpers/arrays"
	"github.com/stretchr/testify/assert"
	"github.com/ugurcsen/gods-generic/lists/arraylist"
)

func TestIsEmpty(t *testing.T) {
	var values []int

	actual := arrays.IsEmpty(values)
	assert.True(t, actual)
}

func TestIsNotEmpty(t *testing.T) {
	values := arraylist.New[int]()
	values.Add(1, 2)

	assert.False(t, values.Empty())
}

func TestContains(t *testing.T) {
	var values []int
	values = append(values, 1, 2)

	actual := arrays.Contains(values, 1)
	assert.True(t, actual)
}

func TestContains_NotFound(t *testing.T) {
	var values []int
	values = append(values, 1, 2)

	actual := arrays.Contains(values, 3)
	assert.False(t, actual)
}

func TestFind(t *testing.T) {
	var values []int
	values = append(values, 1, 2)

	actual := arrays.
		Find(values, func(n int) bool { return n == 2 })
	assert.NotNil(t, actual)
}

func TestFind_NotFound(t *testing.T) {
	var values []int
	values = append(values, 1, 2)

	actual := arrays.
		Find(values, func(n int) bool { return n == 3 })
	assert.Nil(t, actual)
}

func TestFind_ArrayList(t *testing.T) {
	values := arraylist.New[int]()
	values.Add(1, 2, 3)

	result := values.
		Select(func(index int, value int) bool { return value > 1 })

	assert.Equal(t, 2, result.Size())
}
