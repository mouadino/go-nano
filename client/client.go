/*
Package client represents an RPC client to make request to a remote services.

Example:

		c := client.New("upper", http.New())
		reply, err := c.Call("Upper", "foo")
		fmt.Println(reply)

Using asynchronous api:

		f := c.Go("Upper", "foo")
		reply, err := f.Result()
		fmt.Println(reply)

With discovery:

		zk := zookeeper.New("127.0.0.1:2181")
		c := client.New("upper", http.New(), loadbalancer.New(zk, loadbalancer.NewRoundRobin()))
		reply, err := c.Call("Upper", "foo")
		fmt.Println(reply)

*/
package client

import (
	"github.com/mouadino/go-nano/client/extension"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
	"github.com/mouadino/go-nano/utils"
)

// TODO: Add timeout here.

// Client represents an RPC client.
type Client struct {
	endpoint string
	sender   transport.Sender
}

// NewClient returns a new RPC Client customized with specific protocol and extensions.
func New(endpoint string, sender transport.Sender, exts ...extension.Extension) Client {
	return Client{
		endpoint: endpoint,
		sender:   extension.Decorate(sender, exts...),
	}
}

// TODO: Implement transport.Sender ?
// Send a raw request to endpoint, waits service to respond.
func (c *Client) Send(req *protocol.Request) (*protocol.Response, error) {
	return c.sender.Send(c.endpoint, req)
}

// Call a remote method with given parameters, waits service to respond.
func (c *Client) Call(method string, params ...interface{}) (interface{}, error) {
	req := protocol.Request{
		Method: method,
		Params: utils.ParamsFormat(params...),
		// TODO: headers ?
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

// Go calls a remote function asynchronously and returns a future.
func (c *Client) Go(method string, params ...interface{}) *Future {
	f := newFuture()
	go func() {
		resp, err := c.Call(method, params...)
		f.set(resp, err)
	}()
	return f
}
