package client

import (
	"testing"

	"github.com/mouadino/go-nano/protocol"
)

type dummySender struct{}

func (dummySender) Send(endpoint string, req *protocol.Request) (*protocol.Response, error) {
	return &protocol.Response{
		Body: "foobar",
	}, nil
}

func TestClientCall(t *testing.T) {
	c := New("<memory>", dummySender{})

	resp, err := c.Call("Foo", "Arg1", "Arg2")

	if err != nil {
		t.Errorf("Didn't expect to fail, got %s", err)
	}

	if resp.(string) != "foobar" {
		t.Errorf("c.Call(...) got %q, want %q", resp, "foobar")
	}
}

func TestClientGo(t *testing.T) {
	c := New("<memory>", dummySender{})

	f := c.Go("Foo", "Arg1", "Arg2")

	resp, err := f.Result()

	if err != nil {
		t.Errorf("Didn't expect to fail, got %s", err)
	}

	if resp.(string) != "foobar" {
		t.Errorf("c.Call(...) got %q, want %q", resp, "foobar")
	}
}
