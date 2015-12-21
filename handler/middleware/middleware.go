/*
package middleware contains common handler middlewares, example: logger,
tracer, recoverer ... .

*/
package middleware

import "github.com/mouadino/go-nano/handler"

// Chain creates a Handler with the given middlewares stacked in it.
func Chain(hdlr handler.Handler, middlewares ...handler.Middleware) handler.Handler {
	for _, m := range middlewares {
		hdlr = m(hdlr)
	}
	return hdlr
}
