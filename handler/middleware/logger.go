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

func NewLoggerMiddleware(logger *log.Logger) handler.Middleware {
	return func(wrapped handler.Handler) handler.Handler {
		return &loggerMiddleware{
			logger:  logger,
			wrapped: wrapped,
		}
	}
}

func (m *loggerMiddleware) Handle(rw protocol.ResponseWriter, req *protocol.Request) {
	m.logger.WithFields(log.Fields{
		"method": req.Method,
	}).Info("Calling")
	start := time.Now()

	m.wrapped.Handle(rw, req)

	m.logger.WithFields(log.Fields{
		"method":   req.Method,
		"duration": time.Since(start),
	}).Info("Call")
}
