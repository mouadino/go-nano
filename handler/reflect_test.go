package handler

import (
	"errors"
	"testing"

	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/utils"
)

type EchoService struct{}

func (s *EchoService) NanoStart() error {
	return nil
}

func (s *EchoService) Echo(text string) (string, error) {
	return text, nil
}

func (s *EchoService) Fail(text string) (string, error) {
	return "", errors.New("Fail!")
}

func TestReflection(t *testing.T) {
	handler := newStructHandler(&EchoService{})

	if len(handler.methods) != 2 {
		t.Errorf("Expected %d method, got %d", 2, len(handler.methods))
	}
}

func TestHandling(t *testing.T) {
	handler := newStructHandler(&EchoService{})
	req := protocol.Request{
		Method: "echo",
		Params: utils.ParamsFormat("foobar"),
	}
	resp := &protocol.Response{}

	handler.Handle(resp, &req)

	ret, ok := resp.Body.(string)
	if !ok {
		t.Errorf("Return data is not string")
	}

	if ret != "foobar" {
		t.Errorf("Expected handler to return %s, got %s", "foobar", resp.Body)
	}
}

func TestUnknownMethod(t *testing.T) {
	handler := newStructHandler(&EchoService{})
	req := protocol.Request{
		Method: "blabla",
		Params: utils.ParamsFormat("foobar"),
	}
	resp := &protocol.Response{}

	handler.Handle(resp, &req)

	if resp.Error != protocol.UnknownMethod {
		t.Errorf("Expected handle to fail with %s, got %s", protocol.UnknownMethod, resp.Error)
	}
}

func TestWrongArgumentsMethod(t *testing.T) {
	handler := newStructHandler(&EchoService{})
	req := protocol.Request{
		Method: "echo",
		Params: utils.ParamsFormat("foobar", 1),
	}
	resp := &protocol.Response{}

	handler.Handle(resp, &req)

	if resp.Error != protocol.ParamsError {
		t.Errorf("Expected handle to fail with %s, got %s", protocol.ParamsError, resp.Error)
	}
}

func TestRemoteError(t *testing.T) {
	handler := newStructHandler(&EchoService{})
	req := protocol.Request{
		Method: "fail",
		Params: utils.ParamsFormat("foobar"),
	}
	resp := &protocol.Response{}

	handler.Handle(resp, &req)

	if resp.Error == nil || resp.Error.Error() != "Fail!" {
		t.Errorf("expected to fail with 'Fail!' got %s", resp.Error)
	}
}
