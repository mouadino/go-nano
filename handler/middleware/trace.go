package middleware

import (
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	"github.com/pborman/uuid"
)

const TraceHeader = "X-Trace-Id"

type traceMiddleware struct {
	wrapped handler.Handler
}

// NewTrace returns a middleware that add a tracing header to
// responses to be able to correlate requests and responses.
func NewTrace() handler.Middleware {
	return func(wrapped handler.Handler) handler.Handler {
		return &traceMiddleware{
			wrapped: wrapped,
		}
	}
}

func (m *traceMiddleware) Handle(resp *protocol.Response, req *protocol.Request) {
	m.wrapped.Handle(resp, req)

	traceId := req.Header.Get(TraceHeader)
	if traceId == "" {
		traceId = uuid.New()
	}
	resp.Header.Set(TraceHeader, traceId)
	// TODO: Logging with trace id.
}
