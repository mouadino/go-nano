package client

import (
	"errors"
	"time"
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
	fail    chan struct{}
	reply   chan reply
}

func NewTimeoutExt(timeout time.Duration) ClientExtension {
	return func(c Client) Client {
		return &TimeoutExt{
			client:  c,
			timeout: timeout,
			fail:    make(chan struct{}, 1),
			reply:   make(chan reply, 1),
		}
	}
}

func (e *TimeoutExt) Call(method string, params ...interface{}) (interface{}, error) {
	go func() {
		data, err := e.client.Call(method, params...)
		e.reply <- reply{data, err}
	}()
	go func() {
		time.Sleep(e.timeout)
		e.fail <- struct{}{}
	}()

	select {
	case <-e.fail:
		return nil, TimeOutError
	case res := <-e.reply:
		return res.data, res.err
	}
}
