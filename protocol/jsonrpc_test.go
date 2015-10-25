package protocol

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/mouadino/go-nano/transport"
)

var (
	expectedBody = []byte(`{"id": "0", "jsonrpc": "2.0", "method": "upper", "params": {"text": "foobar"}}`)
)

func equalJSON(first, second []byte) error {
	f, err := toJSON(first)
	if err != nil {
		return err
	}
	s, err := toJSON(second)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(f, s) {
		return fmt.Errorf("Expected %s, got %s", first, second)
	}
	return nil
}

func toJSON(b []byte) (interface{}, error) {
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %s", err)
	}
	return f, err
}

func TestJSONRPCBodyHandler(t *testing.T) {
	proto := NewJSONRPCProtocol(transport.NewInMemoryTransport())
	req := Request{"upper", Params{"text": "foobar"}}

	body, err := proto.getBody(&req)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if err := equalJSON(body, expectedBody); err != nil {
		t.Errorf("%s", err)
	}
}

func TestReceiveRequest(t *testing.T) {
	trans := transport.NewInMemoryTransport(
		transport.Data{Body: expectedBody, Resp: &transport.DummyResponseWriter{}},
	)
	proto := NewJSONRPCProtocol(trans)

	_, req := proto.ReceiveRequest()

	if req.Method != "upper" {
		t.Errorf("Expected method %s, got %s", "upper", req.Method)
	}

	params := map[string]string{"text": "foobar"}
	if reflect.DeepEqual(req.Params, params) {
		t.Errorf("Expected params %s, got %s", params, req.Params)
	}
}

// TODO: Test for JSONRPCResponseWriter
