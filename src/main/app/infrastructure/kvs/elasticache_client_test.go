package kvs_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/src/main/app/container"
	"github.com/src/main/app/infrastructure/kvs"
	"github.com/src/main/app/model"
	"github.com/stretchr/testify/assert"
)

func TestElasticCacheClient_Save(t *testing.T) {
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](container.ProvideKVSClient())
	statusDTO := new(model.AppStatusDTO)
	statusDTO.Status = model.Started

	err := kvsClient.Save("go-consumer-app", statusDTO)
	assert.NoError(t, err)
}

func TestElasticCacheClient_SaveErr(t *testing.T) {
	kvsClient := kvs.NewElasticCacheClient[string](container.ProvideKVSClient())
	err := kvsClient.Save("key", nil)
	assert.Error(t, err)

	err = kvsClient.Save("", aws.String("value"))
	assert.Error(t, err)

	kvsClient = kvs.NewElasticCacheClient[string](redis.NewClient(&redis.Options{
		Addr: "invalid",
	}))

	err = kvsClient.Save("key", aws.String("value"))
	assert.Error(t, err)

	err = kvsClient.Save("key", nil)
	assert.Error(t, err)
}

func TestElasticCacheClient_Get(t *testing.T) {
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](container.ProvideKVSClient())
	statusDTO := new(model.AppStatusDTO)
	statusDTO.Status = model.Started

	err := kvsClient.Save("go-consumer-app", statusDTO)
	assert.NoError(t, err)

	actual, err := kvsClient.Get("go-consumer-app")
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, model.Started, actual.Status)
}

func TestElasticCacheClient_GetErr(t *testing.T) {
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](container.ProvideKVSClient())
	statusDTO := new(model.AppStatusDTO)
	statusDTO.Status = model.Started

	err := kvsClient.Save("go-consumer-app", statusDTO)
	assert.NoError(t, err)

	kvsBadSerialization := kvs.NewElasticCacheClient[string](container.ProvideKVSClient())
	actual, err := kvsBadSerialization.Get("go-consumer-app")
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestElasticCacheClient_Get_NotFound(t *testing.T) {
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](container.ProvideKVSClient())
	actual, err := kvsClient.Get("missing")

	assert.NoError(t, err)
	assert.Nil(t, actual)
}

func TestNilError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](db)
	mock.ExpectGet("key").SetErr(errors.New("internal error"))

	value, err := kvsClient.Get("key")
	assert.Error(t, err)
	assert.Nil(t, value)
}
