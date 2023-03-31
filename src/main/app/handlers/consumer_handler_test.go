package handlers_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/src/main/app/handlers"
	"github.com/src/main/app/model"
	"github.com/src/main/app/server"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ConsumerHandlerSuite struct {
	suite.Suite
	app             *server.App
	consumerService *MockConsumerService
	consumerHandler handlers.IConsumerHandler
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerHandlerSuite))
}

type MockConsumerService struct {
	mock.Mock
}

func (m *MockConsumerService) GetAppStatus() *model.AppStatusDTO {
	args := m.Called()
	return args.Get(0).(*model.AppStatusDTO)
}

func (m *MockConsumerService) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConsumerService) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (suite *ConsumerHandlerSuite) SetupTest() {
	suite.consumerService = new(MockConsumerService)
	suite.consumerHandler = handlers.NewConsumerHandler(suite.consumerService)
	suite.app = server.New()
	suite.app.Server.Add(http.MethodGet, "/consumer/status", suite.consumerHandler.GetStatus)
	suite.app.Server.Add(http.MethodPut, "/consumer/start", suite.consumerHandler.Start)
	suite.app.Server.Add(http.MethodPut, "/consumer/stop", suite.consumerHandler.Stop)
}

func (suite *ConsumerHandlerSuite) TestConsumerHandler_GetStatus() {
	suite.consumerService.On("GetAppStatus").Return(GetAppStatusDTO())

	request := httptest.NewRequest(http.MethodGet, "/consumer/status", nil)
	response, err := suite.app.Server.Test(request)
	suite.NoError(err)
	suite.NotNil(response)
	suite.Equal(http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	suite.NoError(err)
	suite.NotNil(body)

	suite.Equal("{\"status\":\"started\"}", string(body))
}

func (suite *ConsumerHandlerSuite) TestConsumerHandler_Start() {
	suite.consumerService.On("Start").Return(nil)
	suite.consumerService.On("GetAppStatus").Return(GetAppStatusDTO())

	request := httptest.NewRequest(http.MethodPut, "/consumer/start", nil)
	response, err := suite.app.Server.Test(request)
	suite.NoError(err)
	suite.NotNil(response)
	suite.Equal(http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	suite.NoError(err)
	suite.NotNil(body)

	suite.Equal("{\"status\":\"started\"}", string(body))
}

func (suite *ConsumerHandlerSuite) TestConsumerHandler_Stop() {
	suite.consumerService.On("Stop").Return(nil)
	appStatusDTO := GetAppStatusDTO()
	appStatusDTO.Status = model.Stopped
	suite.consumerService.On("GetAppStatus").Return(appStatusDTO)

	request := httptest.NewRequest(http.MethodPut, "/consumer/stop", nil)
	response, err := suite.app.Server.Test(request)
	suite.NoError(err)
	suite.NotNil(response)
	suite.Equal(http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	suite.NoError(err)
	suite.NotNil(body)

	suite.Equal("{\"status\":\"stopped\"}", string(body))
}

func (suite *ConsumerHandlerSuite) TestConsumerHandler_StartErr() {
	suite.consumerService.On("Start").Return(errors.New("timeout"))
	appStatusDTO := GetAppStatusDTO()
	appStatusDTO.Status = model.Stopped
	suite.consumerService.On("GetAppStatus").Return(appStatusDTO)

	request := httptest.NewRequest(http.MethodPut, "/consumer/start", nil)
	response, err := suite.app.Server.Test(request)
	suite.NoError(err)
	suite.NotNil(response)
	suite.Equal(http.StatusInternalServerError, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	suite.NoError(err)
	suite.NotNil(body)

	suite.Equal("{\"status_code\":500,\"message\":\"timeout\"}", string(body))
}

func (suite *ConsumerHandlerSuite) TestConsumerHandler_StopErr() {
	suite.consumerService.On("Stop").Return(errors.New("timeout"))
	appStatusDTO := GetAppStatusDTO()
	appStatusDTO.Status = model.Stopped
	suite.consumerService.On("GetAppStatus").Return(appStatusDTO)

	request := httptest.NewRequest(http.MethodPut, "/consumer/stop", nil)
	response, err := suite.app.Server.Test(request)
	suite.NoError(err)
	suite.NotNil(response)
	suite.Equal(http.StatusInternalServerError, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	suite.NoError(err)
	suite.NotNil(body)

	suite.Equal("{\"status_code\":500,\"message\":\"timeout\"}", string(body))
}

func GetAppStatusDTO() *model.AppStatusDTO {
	appStatusDTO := new(model.AppStatusDTO)
	appStatusDTO.Status = model.Started
	return appStatusDTO
}
