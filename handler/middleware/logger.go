package middleware

import (
	"time"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"

	log "github.com/Sirupsen/logrus"
)

type loggerMiddleware struct {
	logger  *log.Logger
	wrapped handler.Handler
}

// NewLogger returns a middleware that log requests with timing data.
func NewLogger(logger *log.Logger) handler.Middleware {
	return func(wrapped handler.Handler) handler.Handler {
		return &loggerMiddleware{
			logger:  logger,
			wrapped: wrapped,
		}
	}
}

func (m *loggerMiddleware) Handle(resp *protocol.Response, req *protocol.Request) {
	m.logger.WithFields(log.Fields{
		"method": req.Method,
	}).Info("Calling")
	start := time.Now()

	m.wrapped.Handle(resp, req)

	m.logger.WithFields(log.Fields{
		"method":   req.Method,
		"duration": time.Since(start),
	}).Info("Call")
}
