package jsonrpc

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
)

var (
	dummyRequest           = []byte(`{"id": "0", "jsonrpc": "2.0", "method": "upper", "params": {"text": "foobar"}}`)
	dummyResponse          = []byte(`{"id": "0", "jsonrpc": "2.0", "result": "foobar", "error": null}`)
	dummyResponseWithError = []byte(`{"id": "0", "jsonrpc": "2.0", "result": null, "error": {"code": "-32000", "message": "Server error", "data":"Server error"}}`)
)

func equalJSON(b1, b2 []byte) bool {
	var (
		m1 map[string]interface{}
		m2 map[string]interface{}
	)

	json.Unmarshal(b1, &m1)
	json.Unmarshal(b2, &m2)

	return reflect.DeepEqual(m1, m2)
}

func TestDecodeRequest(t *testing.T) {
	proto := New()
	req, err := proto.DecodeRequest(bytes.NewReader(dummyRequest))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if req.Method != "upper" {
		t.Errorf("Expected method %s, got %s", "upper", req.Method)
	}

	params := map[string]string{"text": "foobar"}
	if reflect.DeepEqual(req.Params, params) {
		t.Errorf("Expected params %s, got %s", params, req.Params)
	}
}

func TestEncodeRequest(t *testing.T) {
	proto := New()

	req := &protocol.Request{
		Method: "upper",
		Params: protocol.Params{"text": "foobar"},
		Header: header.Header{},
	}

	body, err := proto.EncodeRequest(req)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !equalJSON(body, dummyRequest) {
		t.Errorf("EncodeRequest() got %s, want %s", body, dummyRequest)
	}
}

func TestDecodeResponse(t *testing.T) {
	proto := New()
	resp, err := proto.DecodeResponse(bytes.NewReader(dummyResponse))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Body != "foobar" {
		t.Errorf("DecodeResponse() got %s, want %s", resp.Body, "foobar")
	}
}

func TestEncodeResponse(t *testing.T) {
	proto := New()
	resp := &protocol.Response{
		Body: "foobar",
	}
	body, err := proto.EncodeResponse(resp)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !equalJSON(body, dummyResponse) {
		t.Errorf("EncodeResponse() got %s, want %s", body, dummyResponse)
	}
}

func TestDecodeResponseWithError(t *testing.T) {
	proto := New()
	resp, err := proto.DecodeResponse(bytes.NewReader(dummyResponseWithError))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(resp.Error, protocol.ServerError) {
		t.Errorf("resp.Error got %s, want %s", resp.Error, protocol.ServerError)
	}
}

func TestEncodeResponseWithError(t *testing.T) {
	proto := New()
	resp := &protocol.Response{
		Error: protocol.ServerError,
	}
	body, err := proto.EncodeResponse(resp)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !equalJSON(body, dummyResponseWithError) {
		t.Errorf("EncodeResponse() got %s, want %s", body, dummyResponseWithError)
	}
}
