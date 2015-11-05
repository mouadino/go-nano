package nano

import (
	"time"

	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

func DefaultClient(endpoint string) client.Client {
	return CustomClient(
		endpoint,
		protocol.NewJSONRPCProtocol(transport.NewHTTPTransport()),
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
