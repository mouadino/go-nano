package handler

import "github.com/mouadino/go-nano/protocol"

type Handler interface {
	Handle(protocol.ResponseWriter, *protocol.Request)
}

type Middleware func(Handler) Handler
