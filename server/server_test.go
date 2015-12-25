package server

import (
	"testing"
	"time"

	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/dummy"
)

func buildTestServer(rw *dummy.ResponseWriter, req *protocol.Request) *Server {
	proto := dummy.New(rw, req)
	server := New(proto)

	return server
}

func echoHandler(rw protocol.ResponseWriter, req *protocol.Request) {
	rw.Set(req.Params["msg"])
}

type dummyAnnouncer struct {
	instances map[string]discovery.Instance
}

func (an *dummyAnnouncer) Announce(name string, inst discovery.Instance) error {
	an.instances[name] = inst
	return nil
}

func TestMultipleRegister(t *testing.T) {
	rw := dummy.NewResponseRecorder()
	req := &protocol.Request{
		Method: "echo",
		Params: protocol.Params{"msg": "foobar"},
	}
	server := buildTestServer(rw, req)

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
	rw := dummy.NewResponseRecorder()
	req := &protocol.Request{
		Method: "echo",
		Params: protocol.Params{"msg": "foobar"},
	}
	server := buildTestServer(rw, req)

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
	rw := dummy.NewResponseRecorder()
	req := &protocol.Request{
		Method: "echo",
		Params: protocol.Params{"msg": "foobar"},
	}
	server := buildTestServer(rw, req)

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
	rw := dummy.NewResponseRecorder()
	req := &protocol.Request{
		Method: "echo",
		Params: protocol.Params{"msg": "foobar"},
	}
	server := buildTestServer(rw, req)

	err := server.Register("echo", handler.HandlerFunc(echoHandler))
	if err != nil {
		t.Fatalf("unexpected failure %s", err)
	}

	server.Serve()

	time.Sleep(1 * time.Second)

	if rw.Data.(string) != "foobar" {
		t.Errorf("handler want %q, got %q", "foobar", rw.Data)
	}
}

func TestServeUnknownHandler(t *testing.T) {
	// FIXME:
	t.Skip("Fail on go tip")
	rw := dummy.NewResponseRecorder()
	req := &protocol.Request{
		Method: "UnknownMethod",
	}
	server := buildTestServer(rw, req)

	err := server.Register("echo", handler.HandlerFunc(echoHandler))
	if err != nil {
		t.Fatalf("unexpected failure %s", err)
	}

	server.Serve()

	time.Sleep(1 * time.Second)

	if rw.Error == nil || rw.Error.Error() != `Unknown handler "UnknownMethod"` {
		t.Errorf("expected error got %q", rw.Error)
	}
}
