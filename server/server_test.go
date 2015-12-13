package server

import (
	"testing"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/dummy"
)

func echoHandler(rw protocol.ResponseWriter, req *protocol.Request) {
	rw.Set(req.Params["msg"])
}

func TestMuxServer(t *testing.T) {
	proto := dummy.New()
	server := New(proto)

	err := server.Register("echo", handler.HandlerFunc(echoHandler))
	if err != nil {
		t.Errorf("unexpected failure %s", err)
	}

	rw := &dummy.ResponseWriter{}
	req := &protocol.Request{
		Method: "echo",
		Params: protocol.Params{"msg": "foobar"},
	}

	server.mux.Handle(rw, req)

	if rw.Data.(string) != "foobar" {
		t.Errorf("handler want %q, got %q", "foobar", rw.Data)
	}
}

func TestMultipleRegister(t *testing.T) {
	proto := dummy.New()
	server := New(proto)

	err := server.Register("echo", handler.HandlerFunc(echoHandler))
	if err != nil {
		t.Errorf("unexpected failure %s", err)
	}

	err = server.Register("echo", handler.HandlerFunc(echoHandler))
	if err == nil {
		t.Errorf("expected to fail when registering with same name")
	}
}

func TestUnknownHandler(t *testing.T) {
	proto := dummy.New()

	server := New(proto)
	err := server.Register("echo", handler.HandlerFunc(echoHandler))
	if err != nil {
		t.Errorf("unexpected failure %s", err)
	}

	rw := &dummy.ResponseWriter{}
	invalidReq := &protocol.Request{
		Method: "UnknownMethod",
	}

	server.mux.Handle(rw, invalidReq)

	if rw.Error == nil {
		t.Error("expected to fail with unknown handler")
	}
}
