package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/src/main/app/container"
	"github.com/src/main/app/infrastructure/kvs"
	"github.com/src/main/app/model"
	"github.com/src/main/app/services"
	"github.com/stretchr/testify/assert"
)

func TestConsumerService_IsEnabled(t *testing.T) {
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](container.ProvideKVSClient())
	consumerService := services.NewConsumerService(kvsClient)

	actual := consumerService.GetAppStatus()
	assert.Equal(t, model.Started, actual.Status)
}

func TestConsumerService_Start(t *testing.T) {
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](container.ProvideKVSClient())
	consumerService := services.NewConsumerService(kvsClient)

	err := consumerService.Start()
	assert.NoError(t, err)
	err = consumerService.Start()
	assert.NoError(t, err)

	actual := consumerService.GetAppStatus()
	assert.Equal(t, model.Started, actual.Status)
}

func TestConsumerService_StartGetErr(t *testing.T) {
	db, mock := redismock.NewClientMock()
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](db)
	consumerService := services.NewConsumerService(kvsClient)

	err := consumerService.Start()
	mock.ExpectGet("consumers:go-consumer-app:v1").RedisNil()
	assert.Error(t, err)

	actual := consumerService.GetAppStatus()
	assert.Equal(t, model.Started, actual.Status)
}

func TestConsumerService_StartSaveErr(t *testing.T) {
	db, mock := redismock.NewClientMock()
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](db)
	consumerService := services.NewConsumerService(kvsClient)

	appStatusDTO := new(model.AppStatusDTO)
	appStatusDTO.Status = model.Started

	db.Set(context.TODO(), "consumers:go-consumer-app:v1", appStatusDTO, 0)
	mock.ExpectGet("consumers:go-consumer-app:v1").RedisNil()
	mock.ExpectSet("consumers:go-consumer-app:v1", appStatusDTO, 0).SetErr(errors.New("timeout"))

	err := consumerService.Start()
	assert.Error(t, err)

	actual := consumerService.GetAppStatus()
	assert.Equal(t, model.Started, actual.Status)
}

func TestConsumerService_Stop(t *testing.T) {
	kvsClient := kvs.NewElasticCacheClient[model.AppStatusDTO](container.ProvideKVSClient())
	consumerService := services.NewConsumerService(kvsClient)
	err := consumerService.Stop()
	assert.NoError(t, err)

	actual := consumerService.GetAppStatus()
	assert.Equal(t, model.Stopped, actual.Status)
}
