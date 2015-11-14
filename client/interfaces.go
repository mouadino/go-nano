package client

import "github.com/mouadino/go-nano/protocol"

type Client interface {
	CallEndpoint(string, *protocol.Request) (interface{}, error)
}

type ClientExtension func(Client) Client
