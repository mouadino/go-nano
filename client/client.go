package client

import (
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/utils"
)

type DefaultClient struct {
	Endpoint string
	Proto    protocol.Protocol
}

func (c *DefaultClient) Call(method string, params ...interface{}) (interface{}, error) {
	req := protocol.Request{
		Method: method,
		Params: utils.ParamsFormat(params...),
	}
	resp, err := c.Proto.SendRequest(c.Endpoint, &req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Decorate(c Client, exts ...ClientExtension) Client {
	for _, ext := range exts {
		c = ext(c)
	}
	return c
}
