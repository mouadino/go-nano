package middleware

import (
	"testing"

	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/pborman/uuid"
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

	if traceID == "" {
		t.Error("expected a trace header got nothing")
	}
}

func TestTraceMiddlewareDelegation(t *testing.T) {
	traceID := uuid.New()
	req := &protocol.Request{
		Method: "foobar",
		Params: protocol.Params{},
		Header: header.Header{
			TraceHeader: traceID,
		},
	}
	rw := &protocol.DumpResponseWriter{
		HeaderValues: header.Header{},
	}

	handler := Chain(&DumpHandler{}, NewTraceMiddleware())

	handler.Handle(rw, req)

	newTraceID := rw.Header().Get(TraceHeader)
	if newTraceID != traceID {
		t.Errorf("trace id didn't match, expected %s, got %s", traceID, newTraceID)
	}
}
