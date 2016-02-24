package middleware

import "github.com/mouadino/go-nano/protocol"

type dummyHandler struct{}

func (h *dummyHandler) Handle(resp *protocol.Response, req *protocol.Request) {
}
