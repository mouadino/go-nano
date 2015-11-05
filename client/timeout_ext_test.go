package client

import (
	"testing"
	"time"
)

type DummyClient struct {
	interval time.Duration
}

func (c *DummyClient) Call(m string, p ...interface{}) (interface{}, error) {
	time.Sleep(c.interval)
	return nil, nil
}

func TestTimeoutTrigger(t *testing.T) {
	c := Decorate(
		&DummyClient{2 * time.Second},
		NewTimeoutExt(1*time.Second),
	)

	_, err := c.Call("foobar")

	if err != TimeOutError {
		t.Error("Timeout didn't get trigger else it should")
	}
}

func TestNoTimeout(t *testing.T) {
	c := Decorate(
		&DummyClient{1 * time.Second},
		NewTimeoutExt(5*time.Second),
	)

	_, err := c.Call("foobar")

	if err != nil {
		t.Error("Timeout was trigger else it shouldn't")
	}
}
