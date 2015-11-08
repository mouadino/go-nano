package middleware

import (
	"testing"

	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	uuid "github.com/nu7hatch/gouuid"
)

func TestTraceMiddlewareGeneration(t *testing.T) {
	req := &protocol.Request{
		Method: "foobar",
		Params: protocol.Params{},
		Header: header.Header{},
	}
	rw := &protocol.DumpResponseWriter{
		HeaderValues: header.Header{},
	}

	handler := Chain(&DumpHandler{}, NewTraceMiddleware())

	handler.Handle(rw, req)

	traceID := rw.Header().Get(TraceHeader)

	_, err := uuid.ParseHex(traceID)
	if err != nil {
		t.Errorf("Invalid trace id %s: %s", traceID, err)
	}
}

func TestTraceMiddlewareDelegation(t *testing.T) {
	traceID, _ := uuid.NewV4()
	req := &protocol.Request{
		Method: "foobar",
		Params: protocol.Params{},
		Header: header.Header{
			TraceHeader: traceID.String(),
		},
	}
	rw := &protocol.DumpResponseWriter{
		HeaderValues: header.Header{},
	}

	handler := Chain(&DumpHandler{}, NewTraceMiddleware())

	handler.Handle(rw, req)

	newTraceID := rw.Header().Get(TraceHeader)
	if newTraceID != traceID.String() {
		t.Errorf("trace id didn't match, expected %s, got %s", traceID, newTraceID)
	}
}
