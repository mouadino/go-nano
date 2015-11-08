package middleware

import "github.com/mouadino/go-nano/handler"

func Chain(handler handler.Handler, middlewares ...handler.Middleware) handler.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}
