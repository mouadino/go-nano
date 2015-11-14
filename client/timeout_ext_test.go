package client

import (
	"testing"
	"time"

	"github.com/mouadino/go-nano/protocol"
)

type DummyClient struct {
	interval time.Duration
}

func (c *DummyClient) CallEndpoint(m string, r *protocol.Request) (interface{}, error) {
	// Simulate a real RPC.
	time.Sleep(10 * time.Millisecond)
	return nil, nil
}

func TestTimeoutTrigger(t *testing.T) {
	c := Decorate(
		&DummyClient{},
		NewTimeoutExt(5*time.Millisecond),
	)

	_, err := c.CallEndpoint("foobar", &protocol.Request{})

	if err != TimeOutError {
		t.Error("Timeout didn't get trigger else it should")
	}
}

func TestNoTimeout(t *testing.T) {
	c := Decorate(
		&DummyClient{},
		NewTimeoutExt(1*time.Second),
	)

	_, err := c.CallEndpoint("foobar", &protocol.Request{})

	if err != nil {
		t.Error("Timeout was trigger else it shouldn't")
	}
}
