/*
Package client represents an RPC client that is enable to make request
to remote services.

Example:

		client := DefaultClient("http://127.0.0.1:8080")
		reply, err := client.Call("Upper", "foo")
		fmt.Println(reply)

Using asynchronous api:

		client := DefaultClient("http://127.0.0.1:8080")
		f := client.Go("Upper", "foo")
    // Other code ...
		reply, err := f.Result()
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

// Future represents the response of an asynchronous client request.
type Future struct {
	resp     interface{}
	err      error
	finish   chan struct{}
	finished bool
}

func newFuture() *Future {
	return &Future{
		finish: make(chan struct{}, 1),
	}
}

func (f *Future) Result() (interface{}, error) {
	if !f.finished {
		<-f.finish
	}
	return f.resp, f.err
}

func (f *Future) set(resp interface{}, err error) {
	f.resp = resp
	f.err = err
	f.finished = true
	f.finish <- struct{}{}
}

// Client represents an RPC client.
type Client struct {
	endpoint string
	sender   protocol.Sender
}

// DefaultClient returns a new nano.Client using default configuration.
// Default configuration include default 3 seconds timeout, discovery using
// zookeeper with round robin load balancing strategy.
func DefaultClient(endpoint string) Client {
	// TODO: Protocol factory from endpoint scheme.
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

// Go calls a remote function asynchronously and returns a future object.
func (c *Client) Go(method string, params ...interface{}) *Future {
	f := newFuture()
	go func() {
		resp, err := c.Call(method, params...)
		f.set(resp, err)
	}()
	return f
}
