package middleware

import "github.com/mouadino/go-nano/protocol"

type DumpHandler struct{}

func (h *DumpHandler) Handle(rw protocol.ResponseWriter, req *protocol.Request) {
}
