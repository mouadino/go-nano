package extension

import (
	"errors"
	"time"

	"github.com/mouadino/go-nano/protocol"
)

var (
	TimeOutError = errors.New("Timeout")
)

type reply struct {
	resp *protocol.Response
	err  error
}

type timeoutExt struct {
	sender  protocol.Sender
	timeout time.Duration
}

// NewTimeoutExt returns an extension that wraps a client to timeout
// a request when this later take more than given duration.
func NewTimeoutExt(timeout time.Duration) Extension {
	return func(s protocol.Sender) protocol.Sender {
		return &timeoutExt{
			sender:  s,
			timeout: timeout,
		}
	}
}

func (e *timeoutExt) Send(endpoint string, req *protocol.Request) (*protocol.Response, error) {
	res := make(chan reply, 1)

	go func() {
		resp, err := e.sender.Send(endpoint, req)
		res <- reply{resp, err}
	}()

	select {
	case <-time.After(e.timeout):
		// FIXME: Returing error in second part will break retry strategy.
		return nil, TimeOutError
	case r := <-res:
		return r.resp, r.err
	}
}
