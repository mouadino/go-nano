package nano

import (
	"time"

	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

func Client(endpoint string) *client.Client {
	c := CustomClient(endpoint, protocol.NewJSONRPCProtocol(transport.NewHTTPTransport()))
	c.With(client.NewTimeoutFilter(3 * time.Second))
	return c
}

func CustomClient(endpoint string, proto protocol.Protocol) *client.Client {
	return &client.Client{
		Endpoint: endpoint,
		Proto:    proto,
	}
}
