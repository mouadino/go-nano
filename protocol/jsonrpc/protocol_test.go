package jsonrpc

import (
	"reflect"
	"testing"

	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
)

var (
	dummyRequest           = []byte(`{"id": "0", "jsonrpc": "2.0", "method": "upper", "params": {"text": "foobar"}}`)
	dummyResponse          = []byte(`{"id": "0", "jsonrpc": "2.0", "result": "foobar", "error": null}`)
	dummyResponseWithError = []byte(`{"id": "0", "jsonrpc": "2.0", "result": null, "error": {"code": "-32000", "message": "Server Error"}}`)
)

func TestReceiveRequest(t *testing.T) {
	trans := transport.NewInMemoryTransport(
		[][]byte{dummyRequest}, [][]byte{dummyResponse},
	)
	proto := NewJSONRPCProtocol(trans, serializer.JSONSerializer{})

	_, req := proto.ReceiveRequest()

	if req.Method != "upper" {
		t.Errorf("Expected method %s, got %s", "upper", req.Method)
	}

	params := map[string]string{"text": "foobar"}
	if reflect.DeepEqual(req.Params, params) {
		t.Errorf("Expected params %s, got %s", params, req.Params)
	}
}

func TestSendRequest(t *testing.T) {
	trans := transport.NewInMemoryTransport(
		[][]byte{dummyRequest}, [][]byte{dummyResponse},
	)
	proto := NewJSONRPCProtocol(trans, serializer.JSONSerializer{})

	req := protocol.Request{"upper", protocol.Params{"text": "foobar"}}
	resp, err := proto.SendRequest("", &req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp != "foobar" {
		t.Errorf("Expected response to be 'foobar', was %v", resp)
	}
}

func TestSendRequestWithError(t *testing.T) {
	trans := transport.NewInMemoryTransport(
		[][]byte{dummyRequest}, [][]byte{dummyResponseWithError},
	)
	proto := NewJSONRPCProtocol(trans, serializer.JSONSerializer{})

	req := protocol.Request{"upper", protocol.Params{"text": "foobar"}}
	_, err := proto.SendRequest("", &req)

	if err != protocol.ServerError {
		t.Errorf("Error expected ServerError, else it was %s", err)
	}
}
