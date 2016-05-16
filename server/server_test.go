package server

import (
	"testing"

	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/transport/memory"
)

func echoHandler(resp *protocol.Response, req *protocol.Request) {
	resp.Body = req.Params["msg"]
}

type dummyAnnouncer struct {
	instances map[string]discovery.Instance
}

func (an *dummyAnnouncer) Announce(name string, inst discovery.Instance) error {
	an.instances[name] = inst
	return nil
}

func TestMultipleRegister(t *testing.T) {
	server := New(memory.New(), jsonrpc.New())

	err := server.Register("echo", handler.HandlerFunc(echoHandler))
	if err != nil {
		t.Errorf("unexpected failure %s", err)
	}

	if _, ok := server.mux.hdlrs["echo"]; !ok {
		t.Errorf("'echo' not found in mux handlers")
	}

	err = server.Register("echo", handler.HandlerFunc(echoHandler))
	if err == nil {
		t.Errorf("expected to fail when registering with same name")
	}
}

func TestRegisterWithMetadata(t *testing.T) {
	server := New(memory.New(), jsonrpc.New())

	err := server.RegisterWithMetadata(
		"echo",
		handler.HandlerFunc(echoHandler),
		map[string]interface{}{"datacenter": "eu"},
	)

	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if _, ok := server.mux.hdlrs["echo"]; !ok {
		t.Errorf("'echo' not found in mux handlers")
	}

	if _, ok := server.metas["echo"]; !ok {
		t.Errorf("'echo' not found in metas")
	}
}

func TestAnnounce(t *testing.T) {
	server := New(memory.New(), jsonrpc.New())

	err := server.Register(
		"echo",
		handler.HandlerFunc(echoHandler),
	)

	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	err = server.RegisterWithMetadata(
		"demo",
		handler.HandlerFunc(echoHandler),
		map[string]interface{}{"datacenter": "eu"},
	)

	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	an := &dummyAnnouncer{}
	server.ServeAndAnnounce(an)

	if _, ok := an.instances["echo"]; ok {
		t.Errorf("echo not announced correctly")
	}

	if inst, ok := an.instances["demo"]; ok {
		t.Errorf("demo not announced correctly")

		if inst.Meta["datacenter"] != "eu" {
			t.Errorf("instance metadata expected to contain 'datacenter' else was %s", inst.Meta["datacenter"])
		}
	}
}

func TestServeHandler(t *testing.T) {
	trans := memory.New()
	server := New(trans, jsonrpc.New())

	err := server.Register("echo", handler.HandlerFunc(echoHandler))
	if err != nil {
		t.Fatalf("unexpected failure %s", err)
	}

	server.listen()
	server.serve()

	req := &protocol.Request{
		Method: "echo",
		Params: protocol.Params{"msg": "foobar"},
	}
	resp, _ := trans.Send(":memory:", req)

	if resp.Body.(string) != "foobar" {
		t.Errorf("handler want %q, got %q", "foobar", resp.Body)
	}
}

func TestServeUnknownHandler(t *testing.T) {
	trans := memory.New()
	server := New(trans, jsonrpc.New())

	err := server.Register("echo", handler.HandlerFunc(echoHandler))
	if err != nil {
		t.Fatalf("unexpected failure %s", err)
	}

	server.listen()
	server.serve()

	req := &protocol.Request{
		Method: "UnknownMethod",
	}

	resp, _ := trans.Send(":memory:", req)

	if resp.Error == nil || resp.Error.Error() != `Unknown handler "UnknownMethod"` {
		t.Errorf("want unknown handler error got %q", resp.Error)
	}
}
