package clients

import (
	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/config"
	"net/http"
)

type ITargetAppClient interface {
	PostMessage(targetAppRequest *RequestBody) error
}

type Client struct {
	rb      *rest.RequestBuilder
	baseURL string
}

func NewClient(rb *rest.RequestBuilder) Client {
	return Client{
		rb:      rb,
		baseURL: config.String("target-app.endpoint"),
	}
}

func (c Client) PostMessage(requestBody *RequestBody) error {
	response := c.rb.Post(c.baseURL, requestBody)

	if response.Err != nil {
		return response.Err
	}

	if response.StatusCode != http.StatusOK {
		return fiber.NewError(response.StatusCode, response.String())
	}

	return nil
}

type RequestBody struct {
	Msg string `json:"msg,omitempty"`
}
