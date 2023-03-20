package config_test

import (
	"testing"

	"github.com/src/main/app/config"
	"github.com/stretchr/testify/assert"
)

func TestGetProperty(t *testing.T) {
	actual := config.String("app.name")
	assert.Equal(t, "go-consumer-app", actual)
}

func TestGetProperty_Err(t *testing.T) {
	actual := config.String("missing")
	assert.Equal(t, "", actual)
}

func TestGetIntProperty_Err(t *testing.T) {
	actual := config.Int("missing")
	assert.Equal(t, 0, actual)
}

func TestGetTryIntProperty_Err(t *testing.T) {
	actual := config.TryInt("missing", 1)
	assert.Equal(t, 1, actual)
}
