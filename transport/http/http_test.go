package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
)

type echoHandler struct{}

func (echoHandler) Handle(resp *protocol.Response, req *protocol.Request) {
	resp.Body = req.Params["_0"]
}

func TestHTTPReceive(t *testing.T) {
	trans := New()
	trans.Listen()
	trans.AddHandler(jsonrpc.New(), echoHandler{})
	go trans.Serve()

	body := `{"id": "0", "method": "", "params": {"_0": "world"}}`
	result := `{"jsonrpc":"2.0","result":"world","error":null,"id":"0"}`

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/rpc/", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	trans.handle(resp, req)

	b := resp.Body.String()
	if b != result {
		t.Errorf("request body doesn't match want %v, got %v", result, b)
	}
}

func TestHTTPSend(t *testing.T) {
	trans := New()
	// FIXME: Broken api, nil pointer if this 2 are not called.
	trans.AddHandler(jsonrpc.New(), echoHandler{})
	trans.Listen()
	go trans.Serve()

	req := &protocol.Request{
		Method: "",
		Params: protocol.Params{"_0": "foobar"},
		Header: header.Header{},
	}

	resp, err := trans.Send(trans.Addr(), req)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if resp.Body != "foobar" {
		t.Errorf("unexpected response want %q, got %v", "foobar", resp)
	}
}

func TestHTTPAddr(t *testing.T) {
	trans := New()
	trans.Listen()

	addr, err := url.Parse(trans.Addr())

	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if addr.Scheme != "http" {
		t.Errorf("unexpected scheme want %q, got %q", "http", addr.Scheme)
	}
}

// TODO: Add benchmarks.
