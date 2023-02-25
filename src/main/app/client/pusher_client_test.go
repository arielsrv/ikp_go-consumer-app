package client_test

import (
	"net/http"
	"testing"

	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/src/main/app/client"
	"github.com/src/main/app/server/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRequestBuilder struct {
	mock.Mock
}

func (m *MockRequestBuilder) Get(string) *rest.Response {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) Post(string, interface{}) *rest.Response {
	args := m.Called()
	return args.Get(0).(*rest.Response)
}

func (m *MockRequestBuilder) Put(string, interface{}) *rest.Response {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) Patch(string, interface{}) *rest.Response {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) Delete(string) *rest.Response {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) Head(string) *rest.Response {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) Options(string) *rest.Response {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) AsyncGet(string, func(*rest.Response)) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) AsyncPost(string, interface{}, func(*rest.Response)) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) AsyncPut(string, interface{}, func(*rest.Response)) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) AsyncPatch(string, interface{}, func(*rest.Response)) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) AsyncDelete(string, func(*rest.Response)) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) AsyncHead(string, func(*rest.Response)) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) AsyncOptions(string, func(*rest.Response)) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRequestBuilder) ForkJoin(func(*rest.Concurrent)) {
	//TODO implement me
	panic("implement me")
}

func TestNewHTTPPusherClient(t *testing.T) {
	rb := new(MockRequestBuilder)
	rb.On("Post").Return(getResponse())

	httpPusherClient := client.NewHTTPPusherClient(rb, "https://my.app/news")
	requestBody := new(client.RequestBody)
	requestBody.ID = "1"
	requestBody.Msg = "Hello world"

	err := httpPusherClient.PostMessage(requestBody)
	assert.NoError(t, err)
}

func TestNewHTTPPusherClientErr_400(t *testing.T) {
	rb := new(MockRequestBuilder)
	rb.On("Post").Return(getHTTPErrorResponse(http.StatusBadRequest))

	httpPusherClient := client.NewHTTPPusherClient(rb, "https://my.app/news")
	requestBody := new(client.RequestBody)
	requestBody.ID = "1"
	requestBody.Msg = "Hello world"

	err := httpPusherClient.PostMessage(requestBody)
	assert.Error(t, err)
}

func TestNewHTTPPusherClientErr_500(t *testing.T) {
	rb := new(MockRequestBuilder)
	rb.On("Post").Return(getHTTPErrorResponse(http.StatusInternalServerError))

	httpPusherClient := client.NewHTTPPusherClient(rb, "https://my.app/news")
	requestBody := new(client.RequestBody)
	requestBody.ID = "1"
	requestBody.Msg = "Hello world"

	err := httpPusherClient.PostMessage(requestBody)
	assert.Error(t, err)
}

func TestNewHTTPPusherClientTransportError(t *testing.T) {
	rb := new(MockRequestBuilder)
	rb.On("Post").Return(getErrorResponse())

	httpPusherClient := client.NewHTTPPusherClient(rb, "https://my.app/news")
	requestBody := new(client.RequestBody)
	requestBody.ID = "1"
	requestBody.Msg = "Hello world"

	err := httpPusherClient.PostMessage(requestBody)
	assert.Error(t, err)
}

func getHTTPErrorResponse(statusCode int) *rest.Response {
	response := new(rest.Response)
	response.Response = new(http.Response)
	response.StatusCode = statusCode

	return response
}

func getErrorResponse() *rest.Response {
	response := new(rest.Response)
	response.Response = new(http.Response)
	response.Err = errors.New("connection refused")

	return response
}

func getResponse() *rest.Response {
	response := new(rest.Response)
	response.Response = new(http.Response)
	response.StatusCode = 200

	return response
}
