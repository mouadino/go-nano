package client

import "github.com/mouadino/go-nano/protocol"

type Client interface {
	CallEndpoint(string, *protocol.Request) (interface{}, error)
}

type ClientExtension func(Client) Client

type DefaultClient struct {
	Proto protocol.Protocol
}

func (c *DefaultClient) CallEndpoint(endpoint string, req *protocol.Request) (interface{}, error) {
	// TODO: Protocol factory from endpoint scheme.
	resp, err := c.Proto.SendRequest(endpoint, req)
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
