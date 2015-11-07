package client

import (
	"testing"
	"time"
)

type DummyClient struct {
	interval time.Duration
}

func (c *DummyClient) Call(m string, p ...interface{}) (interface{}, error) {
	// Simulate a real RPC.
	time.Sleep(10 * time.Millisecond)
	return nil, nil
}

func TestTimeoutTrigger(t *testing.T) {
	c := Decorate(
		&DummyClient{},
		NewTimeoutExt(5*time.Millisecond),
	)

	_, err := c.Call("foobar")

	if err != TimeOutError {
		t.Error("Timeout didn't get trigger else it should")
	}
}

func TestNoTimeout(t *testing.T) {
	c := Decorate(
		&DummyClient{},
		NewTimeoutExt(1*time.Second),
	)

	_, err := c.Call("foobar")

	if err != nil {
		t.Error("Timeout was trigger else it shouldn't")
	}
}
