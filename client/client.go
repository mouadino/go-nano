/*
Package client represents RPC client.

Example:

		client := DefaultClient("http://127.0.0.1:8080")
		reply, err := client.Call("Upper", "foo")
		fmt.Println(reply)

*/
package client

import (
	"time"

	"github.com/mouadino/go-nano/client/extension"
	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/discovery/loadbalancer"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
	"github.com/mouadino/go-nano/utils"
)

// Client represent an RPC client.
type Client struct {
	endpoint string
	sender   protocol.Sender
}

// DefaultClient returns a new nano.Client using default configuration.
// Default configuration include default 3 seconds timeout, discovery using
// zookeeper with round robin load balancing strategy.
func DefaultClient(endpoint string) Client {
	zkDiscover := discovery.DefaultZooKeeperAnnounceResolver(
		[]string{"127.0.0.1:2181"},
	)
	return CustomClient(
		endpoint,
		jsonrpc.NewJSONRPCProtocol(transport.NewHTTPTransport(), serializer.JSONSerializer{}),
		extension.NewTimeoutExt(3*time.Second),
		loadbalancer.NewLoadBalancerExtension(
			zkDiscover,
			loadbalancer.RoundRobinLoadBalancer(),
		),
	)
}

// CustomClient returns a new nano.Client customized with specific protocol and extensions.
func CustomClient(endpoint string, sender protocol.Sender, exts ...extension.Extension) Client {
	return Client{
		endpoint: endpoint,
		sender:   extension.Decorate(sender, exts...),
	}
}

// Send a raw protocol request to client endpoint, waits service to respond, and returns
// either an error or service response.
func (c *Client) Send(req *protocol.Request) (*protocol.Response, error) {
	// TODO: Protocol factory from endpoint scheme.
	return c.sender.Send(c.endpoint, req)
}

// Call a remote method with given parameters, waits service to respond, and returns
// either an error or service response.
func (c *Client) Call(method string, params ...interface{}) (interface{}, error) {
	req := protocol.Request{
		Method: method,
		Params: utils.ParamsFormat(params...),
	}
	resp, err := c.Send(&req)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Body, nil
}

// Go calls a remote function asynchronously.
func (c *Client) Go(method string, params ...interface{}) {
	// TODO
}
