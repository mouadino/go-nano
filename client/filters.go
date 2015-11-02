package client

import (
	"errors"
	"time"

	"github.com/mouadino/go-nano/protocol"
)

var (
	TimeOutError = errors.New("Timeout")
)

type Res struct {
	Data interface{}
	Err  error
}

type TimeoutFilter struct {
	client  IClient
	timeout time.Duration
	fail    chan struct{}
	finish  chan Res
}

func NewTimeoutFilter(timeout time.Duration) Filter {
	return func(client IClient) IClient {
		return &TimeoutFilter{
			client:  client,
			timeout: timeout,
			fail:    make(chan struct{}, 1),
			finish:  make(chan Res, 1),
		}
	}
}

func (f *TimeoutFilter) SendRequest(endpoint string, req *protocol.Request) (interface{}, error) {
	go func() {
		data, err := f.client.SendRequest(endpoint, req)
		f.finish <- Res{data, err}
	}()
	go func() {
		time.Sleep(f.timeout)
		f.fail <- struct{}{}
	}()

	select {
	case <-f.fail:
		return nil, TimeOutError
	case res := <-f.finish:
		return res.Data, res.Err
	}
}

// TODO: type DiscoveryFilter struct {}
