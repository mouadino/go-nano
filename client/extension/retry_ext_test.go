package extension

import (
	"errors"
	"testing"

	"github.com/cenkalti/backoff"
	"github.com/mouadino/go-nano/protocol"
)

type dummyRetrySender struct {
	called int32
}

func (c *dummyRetrySender) Send(m string, r *protocol.Request) (*protocol.Response, error) {
	c.called++
	return nil, errors.New("FAIL!")
}

func TestRetryTrigger(t *testing.T) {
	sender := &dummyRetrySender{}
	c := NewRetryExt(3, &backoff.ZeroBackOff{})(sender)

	_, err := c.Send("foobar", &protocol.Request{})

	if err != RetryExceeded {
		t.Errorf("retrying sending request want %s, got %s", RetryExceeded, err)
	}

	if sender.called != 3 {
		t.Errorf("sender.called want %d, got %d", 3, sender.called)
	}
}
