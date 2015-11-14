package client

import (
	"errors"
	"time"

	"github.com/mouadino/go-nano/protocol"
)

var (
	TimeOutError = errors.New("Timeout")
)

type reply struct {
	data interface{}
	err  error
}

type TimeoutExt struct {
	client  Client
	timeout time.Duration
	reply   chan reply
}

func NewTimeoutExt(timeout time.Duration) ClientExtension {
	return func(c Client) Client {
		return &TimeoutExt{
			client:  c,
			timeout: timeout,
			reply:   make(chan reply, 1),
		}
	}
}

func (e *TimeoutExt) CallEndpoint(endpoint string, req *protocol.Request) (interface{}, error) {
	go func() {
		data, err := e.client.CallEndpoint(endpoint, req)
		e.reply <- reply{data, err}
	}()

	select {
	case <-time.After(e.timeout):
		return nil, TimeOutError
	case res := <-e.reply:
		return res.data, res.err
	}
}
