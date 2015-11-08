package middleware

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	uuid "github.com/nu7hatch/gouuid"
)

const TraceHeader = "X-Trace-Id"

type traceMiddleware struct {
	wrapped handler.Handler
}

func NewTraceMiddleware() handler.Middleware {
	return func(wrapped handler.Handler) handler.Handler {
		return &traceMiddleware{
			wrapped: wrapped,
		}
	}
}

func (m *traceMiddleware) Handle(rw protocol.ResponseWriter, req *protocol.Request) {
	m.wrapped.Handle(rw, req)

	traceId := req.Header.Get(TraceHeader)

	var err error
	if traceId == "" {
		traceId, err = m.generateUUID()
		if err != nil {
			log.Error("Failed to generate uuid")
			return
		}
	}
	rw.Header().Set(TraceHeader, traceId)
	// TODO: Logging with trace id.
}

func (m *traceMiddleware) generateUUID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
