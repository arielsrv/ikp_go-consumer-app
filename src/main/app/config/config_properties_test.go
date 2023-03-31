package config_test

import (
	"testing"

	"github.com/src/main/app/config"
	"github.com/stretchr/testify/assert"
)

func TestAppConfig(t *testing.T) {
	err := config.MockConfig("config_properties.yml")

	assert.NoError(t, err)

	stringValue := config.String("key")
	assert.Equal(t, "value", stringValue)

	stringValue = config.String("missing")
	assert.Equal(t, "", stringValue)

	boolValue := config.TryBool("enable", true)
	assert.True(t, boolValue)

	boolValue = config.TryBool("logger", true)
	assert.False(t, boolValue)

	intValue := config.TryInt("missing threads", 1)
	assert.Equal(t, 1, intValue)

	intValue = config.TryInt("threads", 1)
	assert.Equal(t, 5, intValue)
}
