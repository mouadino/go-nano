package extension

import (
	"errors"
	"testing"

	"github.com/mouadino/go-nano/protocol"
	"github.com/rubyist/circuitbreaker"
)

type dummyCBSender struct {
	called int32
}

func (c *dummyCBSender) Send(m string, r *protocol.Request) (*protocol.Response, error) {
	c.called++
	return nil, errors.New("FAIL!")
}

func TestCircuitBreakerExtension(t *testing.T) {
	sender := &dummyCBSender{}
	c := NewCircuitBreakerExt(circuit.NewThresholdBreaker(10))(sender)

	_, err := c.Send("foobar", &protocol.Request{})

	if err != OpenCircuitError {
		t.Errorf("circuit breaker wasn't triggered want %s, got %s", OpenCircuitError, err)
	}

	if sender.called != 10 {
		t.Errorf("sender.called want %d, got %d", 10, sender.called)
	}

}
