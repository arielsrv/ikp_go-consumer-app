package services

import (
	"fmt"

	"github.com/src/main/app/config"
	"github.com/src/main/app/infrastructure/kvs"
	"github.com/src/main/app/log"
	"github.com/src/main/app/model"
)

type IConsumerService interface {
	GetAppStatus() *model.AppStatusDTO
	Stop() error
	Start() error
}

type ConsumerService struct {
	kvsClient kvs.Client[model.AppStatusDTO]
}

func NewConsumerService(kvsClient kvs.Client[model.AppStatusDTO]) *ConsumerService {
	return &ConsumerService{kvsClient: kvsClient}
}

func (c ConsumerService) GetAppStatus() *model.AppStatusDTO {
	appStatusDTO := new(model.AppStatusDTO)
	appStatusDTO.Status = model.Started

	cacheKey := getCacheKey()
	resultFromCache, err := c.kvsClient.Get(cacheKey)
	if err != nil {
		log.Warnf("failed to retrieve status from key-value store: %s, started by default", err)
		return appStatusDTO
	}

	if resultFromCache != nil {
		return resultFromCache
	}

	return appStatusDTO
}

func (c ConsumerService) Stop() error {
	return c.refreshKvs(model.Stopped)
}

func (c ConsumerService) Start() error {
	return c.refreshKvs(model.Started)
}

func (c ConsumerService) refreshKvs(status model.Status) error {
	cacheKey := getCacheKey()

	appStatusDTO, err := c.kvsClient.Get(cacheKey)
	if err != nil {
		return err
	}

	if appStatusDTO != nil && appStatusDTO.Status == status {
		log.Warnf("consumer already %s", status)
		return nil
	}

	appStatusDTO = new(model.AppStatusDTO)
	appStatusDTO.Status = status

	err = c.kvsClient.Save(cacheKey, appStatusDTO)
	if err != nil {
		return err
	}

	log.Warnf("consumer switched to %s", status)

	return nil
}

func getCacheKey() string {
	return fmt.Sprintf("consumers:%s:v1", getAppName())
}

func getAppName() string {
	appName := config.String("app.name")
	return appName
}
