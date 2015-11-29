package extension

import (
	"testing"
	"time"

	"github.com/mouadino/go-nano/protocol"
)

type dummySender struct {
	interval time.Duration
}

func (c *dummySender) Send(m string, r *protocol.Request) (*protocol.Response, error) {
	// Simulate a real RPC.
	time.Sleep(10 * time.Millisecond)
	return nil, nil
}

func TestTimeoutTrigger(t *testing.T) {
	c := NewTimeoutExt(5 * time.Millisecond)(&dummySender{})

	_, err := c.Send("foobar", &protocol.Request{})

	if err != TimeOutError {
		t.Error("Timeout didn't get trigger else it should")
	}
}

func TestNoTimeout(t *testing.T) {
	c := NewTimeoutExt(1 * time.Second)(&dummySender{})

	_, err := c.Send("foobar", &protocol.Request{})

	if err != nil {
		t.Error("Timeout was trigger else it shouldn't")
	}
}
