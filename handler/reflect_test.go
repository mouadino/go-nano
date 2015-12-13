package handler

import (
	"testing"

	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/dummy"
	"github.com/mouadino/go-nano/utils"
)

type EchoService struct{}

func (s *EchoService) NanoStart() error {
	return nil
}

func (s *EchoService) Echo(text string) string {
	return text
}

func TestReflection(t *testing.T) {
	handler := NewStructHandler(&EchoService{})

	if len(handler.methods) != 1 {
		t.Errorf("Expected %d method, got %d", 1, len(handler.methods))
	}
}

func TestHandling(t *testing.T) {
	handler := NewStructHandler(&EchoService{})
	req := protocol.Request{
		Method: "echo",
		Params: utils.ParamsFormat("foobar"),
	}
	resp := &dummy.ResponseWriter{}

	handler.Handle(resp, &req)

	ret, ok := resp.Data.(string)
	if !ok {
		t.Errorf("Return data is not string")
	}

	if ret != "foobar" {
		t.Errorf("Expected handler to return %s, got %s", "foobar", resp.Data)
	}
}

func TestUnknownMethod(t *testing.T) {
	handler := NewStructHandler(&EchoService{})
	req := protocol.Request{
		Method: "blabla",
		Params: utils.ParamsFormat("foobar"),
	}
	resp := &dummy.ResponseWriter{}

	handler.Handle(resp, &req)

	if resp.Error != protocol.UnknownMethod {
		t.Errorf("Expected handle to fail with %s, got %s", protocol.UnknownMethod, resp.Error)
	}
}

func TestWrongArgumentsMethod(t *testing.T) {
	handler := NewStructHandler(&EchoService{})
	req := protocol.Request{
		Method: "echo",
		Params: utils.ParamsFormat("foobar", 1),
	}
	resp := &dummy.ResponseWriter{}

	handler.Handle(resp, &req)

	if resp.Error != protocol.ParamsError {
		t.Errorf("Expected handle to fail with %s, got %s", protocol.ParamsError, resp.Error)
	}
}

// TODO: Dealing with headers.
// TODO: RemoteError
