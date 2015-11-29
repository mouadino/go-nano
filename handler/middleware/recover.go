package middleware

import (
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
)

type recoverMiddleware struct {
	wrapped   handler.Handler
	logger    *log.Logger
	showStack bool
	stackSize int
}

func NewRecoverMiddleware(logger *log.Logger, showStack bool, stackSize int) handler.Middleware {
	return func(wrapped handler.Handler) handler.Handler {
		return &recoverMiddleware{
			wrapped:   wrapped,
			logger:    logger,
			showStack: showStack,
			stackSize: stackSize,
		}
	}
}

func (m *recoverMiddleware) Handle(rw protocol.ResponseWriter, req *protocol.Request) {
	defer m.recover(rw)
	m.wrapped.Handle(rw, req)
}

func (m *recoverMiddleware) recover(rw protocol.ResponseWriter) {
	if err := recover(); err != nil {
		rw.SetError(protocol.InternalError)
		m.logger.WithFields(log.Fields{
			"error": err,
		}).Error("Panic")
		if m.showStack {
			m.logger.Error(m.stacktrace())
		}
	}
}

func (m *recoverMiddleware) stacktrace() string {
	stack := make([]byte, m.stackSize)
	return string(stack[:runtime.Stack(stack, false)])
}
