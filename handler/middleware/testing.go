package middleware

import "github.com/mouadino/go-nano/protocol"

type dummyHandler struct{}

func (h *dummyHandler) Handle(rw protocol.ResponseWriter, req *protocol.Request) {
}
