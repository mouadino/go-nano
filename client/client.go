/*
Package client represents an RPC client that is enable to make request
to remote services.

Example:

		c := client.New("upper", jsonrpc.New(http.New()))
		reply, err := c.Call("Upper", "foo")
		fmt.Println(reply)

Using asynchronous api:

		f := c.Go("Upper", "foo")
		reply, err := f.Result()
		fmt.Println(reply)

With discovery:

    zk := zookeeper.New("127.0.0.1:2181")
    c := client.New("upper", jsonrpc.New(http.New()), loadbalancer.New(zk, loadbalancer.NewRoundRobin()))
    reply, err := c.Call("Upper", "foo")
    fmt.Println(reply)

*/
package client

import (
	"github.com/mouadino/go-nano/client/extension"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/utils"
)

// Future represents the response of an asynchronous client request.
type Future struct {
	resp  interface{}
	err   error
	done  chan struct{}
	ready bool
}

func newFuture() *Future {
	return &Future{
		done: make(chan struct{}, 1),
	}
}

// Result returns the future result, block until request finish.
func (f *Future) Result() (interface{}, error) {
	if !f.ready {
		<-f.done
	}
	return f.resp, f.err
}

func (f *Future) set(resp interface{}, err error) {
	f.resp = resp
	f.err = err
	f.ready = true
	f.done <- struct{}{}
}

// Client represents an RPC client.
type Client struct {
	endpoint string
	sender   protocol.Sender
}

// NewClient returns a new RPC Client customized with specific protocol and extensions.
func New(endpoint string, sender protocol.Sender, exts ...extension.Extension) Client {
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
