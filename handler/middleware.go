package handler

import (
	"errors"
	"time"

	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

// TODO: Trace-id Middleware.
// TODO: Log middleware.
// TODO: RateLimit middleware.

var (
	TimeOutError = errors.New("Timeout")
)

type TimeoutMiddleware struct {
	Handler
	timeout time.Duration
	fail    chan struct{}
	finish  chan struct{}
}

func WithTimeout(timeout time.Duration) Middleware {
	return func(h Handler) Handler {
		return &TimeoutMiddleware{
			Handler: h,
			timeout: timeout,
			fail:    make(chan struct{}, 1),
			finish:  make(chan struct{}, 1),
		}
	}
}

func (h *TimeoutMiddleware) Handle(w transport.ResponseWriter, r *protocol.Request) {
	defer close(h.finish)
	defer close(h.fail)
	go func() {
		// TODO: Context and cancellation.
		h.Handler.Handle(w, r)
		h.finish <- struct{}{}
	}()
	go func() {
		time.Sleep(h.timeout * time.Second)
		h.fail <- struct{}{}
	}()

	select {
	case <-h.fail:
		return // FIXME: TimeOutError
	case <-h.finish:
		return
	}
}
