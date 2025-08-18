package client

import (
	"local/model"
)

type Client struct {
	SocketClient SocketClient
}

func NewClient(params *model.InitParams) *Client {
	return &Client{
		SocketClient: NewSocketClient(),
	}
}
