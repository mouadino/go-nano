package client

import (
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/utils"
)

type Client struct {
	Endpoint string
	Proto    IClient
}

func (client *Client) Call(method string, params ...interface{}) (interface{}, error) {
	req := protocol.Request{
		Method: method,
		Params: utils.ParamsFormat(params...),
	}
	resp, err := client.Proto.SendRequest(client.Endpoint, &req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (client *Client) With(filter Filter) {
	client.Proto = filter(client.Proto)
}
