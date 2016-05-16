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

// NewRecover returns a middleware that allow to recover from panic errors.
func NewRecover(logger *log.Logger, showStack bool, stackSize int) handler.Middleware {
	return func(wrapped handler.Handler) handler.Handler {
		return &recoverMiddleware{
			wrapped:   wrapped,
			logger:    logger,
			showStack: showStack,
			stackSize: stackSize,
		}
	}
}

func (m *recoverMiddleware) Handle(resp *protocol.Response, req *protocol.Request) {
	defer m.recover(resp)
	m.wrapped.Handle(resp, req)
}

func (m *recoverMiddleware) recover(resp *protocol.Response) {
	if err := recover(); err != nil {
		resp.Error = protocol.InternalError
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
