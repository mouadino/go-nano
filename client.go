package nano

import (
	"time"

	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
)

func DefaultClient(endpoint string) client.Client {
	return CustomClient(
		endpoint,
		jsonrpc.NewJSONRPCProtocol(transport.NewHTTPTransport(), serializer.JSONSerializer{}),
		client.NewTimeoutExt(3*time.Second),
	)
}

func CustomClient(endpoint string, proto protocol.Protocol, exts ...client.ClientExtension) client.Client {
	c := &client.DefaultClient{
		Endpoint: endpoint,
		Proto:    proto,
	}
	return client.Decorate(c, exts...)
}
