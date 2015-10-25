package handler

import (
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

type Handler interface {
	Handle(transport.ResponseWriter, *protocol.Request)
}

type Middleware func(Handler) Handler
