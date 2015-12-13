package middleware

import (
	log "github.com/Sirupsen/logrus"

	"github.com/mouadino/go-nano/handler"
)

var Defaults = []handler.Middleware{
	NewRecoverMiddleware(log.New(), true, 8*1024),
	NewTraceMiddleware(),
	NewLoggerMiddleware(log.New()),
}

func Chain(hdlr handler.Handler, middlewares ...handler.Middleware) handler.Handler {
	for _, m := range middlewares {
		hdlr = m(hdlr)
	}
	return hdlr
}
