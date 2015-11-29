package extension

// TODO:

import "github.com/mouadino/go-nano/protocol"

type circuitBreakerExt struct {
	sender protocol.Sender
}

func NewCircuitBreakerExt() Extension {
	return func(s protocol.Sender) protocol.Sender {
		return &circuitBreakerExt{
			sender: s,
		}
	}
}

func (e *circuitBreakerExt) Send(endpoint string, req *protocol.Request) (*protocol.Response, error) {
	return e.sender.Send(endpoint, req)
}
