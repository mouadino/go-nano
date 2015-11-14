package nano

import (
	"time"

	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
	"github.com/mouadino/go-nano/utils"
)

type defaultClient struct {
	endpoint string
	client   client.Client
}

func DefaultClient(endpoint string) defaultClient {
	zkDiscover := discovery.DefaultZooKeeperAnnounceResolver(
		[]string{"127.0.0.1:2181"},
	)
	return CustomClient(
		endpoint,
		jsonrpc.NewJSONRPCProtocol(transport.NewHTTPTransport(), serializer.JSONSerializer{}),
		client.NewTimeoutExt(3*time.Second),
		discovery.NewLoadBalancerExtension(
			zkDiscover,
			discovery.RoundRobinLoadBalancer(),
		),
	)
}

func CustomClient(endpoint string, proto protocol.Protocol, exts ...client.ClientExtension) defaultClient {
	c := &client.DefaultClient{
		Proto: proto,
	}
	return defaultClient{
		endpoint: endpoint,
		client:   client.Decorate(c, exts...),
	}
}

func (c *defaultClient) Call(method string, params ...interface{}) (interface{}, error) {
	req := protocol.Request{
		Method: method,
		Params: utils.ParamsFormat(params...),
	}
	return c.client.CallEndpoint(c.endpoint, &req)
}
