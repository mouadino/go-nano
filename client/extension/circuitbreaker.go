package extension

import (
	"errors"

	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
	"github.com/rubyist/circuitbreaker"
)

// OpenCircuitError error is returned when circuit is open.
var OpenCircuitError = errors.New("open circuit breaker")

type circuitBreakerExt struct {
	next transport.Sender
	cb   *circuit.Breaker
}

// NewCircuitBreakerExt returns an extension that implements the circuit breaker
// pattern. Currently no remote exception is counted.
// FIXME: extensions are client specific this mean that all requests share same cb.
func NewCircuitBreakerExt(cb *circuit.Breaker) Extension {
	return func(next transport.Sender) transport.Sender {
		return &circuitBreakerExt{
			next: next,
			cb:   cb,
		}
	}
}

// TODO: Circuit breaker metrics.
// TODO: whitelist of errors.
// TODO: Per instance.
func (e *circuitBreakerExt) Send(endpoint string, req *protocol.Request) (*protocol.Response, error) {
	var (
		resp *protocol.Response
		err  error
	)

	for {
		if e.cb.Ready() {
			resp, err = e.next.Send(endpoint, req)
			if err != nil {
				e.cb.Fail()
				continue
			}
			e.cb.Success()
			return resp, err
		} else {
			return nil, OpenCircuitError
		}
	}
}
