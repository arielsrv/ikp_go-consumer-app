package clients

import (
	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/src/main/app/config"
	"log"
)

type ITargetAppClient interface {
	PostMessage(targetAppRequest RequestBody) error
}

type Client struct {
	rb      *rest.RequestBuilder
	baseURL string
}

func NewClient(rb *rest.RequestBuilder) *Client {
	return &Client{
		rb:      rb,
		baseURL: config.String("rest.client.target-app.baseUrl"),
	}
}

func (c Client) PostMessage(requestBody RequestBody) error {
	log.Println(requestBody)
	//TODO implement me
	panic("implement me")
}

type RequestBody struct {
	Msg string `json:"msg,omitempty"`
}
