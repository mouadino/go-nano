package client

import "github.com/mouadino/go-nano/protocol"

type IClient interface {
	SendRequest(string, *protocol.Request) (interface{}, error)
}

type Filter func(IClient) IClient
