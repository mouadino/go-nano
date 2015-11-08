package handler

import "github.com/mouadino/go-nano/protocol"

type Handler interface {
	Handle(protocol.ResponseWriter, *protocol.Request)
}

type HandlerFunc func(protocol.ResponseWriter, *protocol.Request)

func (h HandlerFunc) Handle(rw protocol.ResponseWriter, req *protocol.Request) {
	h(rw, req)
}

type Middleware func(Handler) Handler
