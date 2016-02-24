package memory

import (
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

type memoryTransport struct {
	proto protocol.Protocol
	hdlr  handler.Handler

	reqs chan *protocol.Request
}

func New() transport.Transport {
	return &memoryTransport{
		reqs: make(chan *protocol.Request),
	}
}

func (trans *memoryTransport) Send(endpoint string, req *protocol.Request) (*protocol.Response, error) {
	resp := &protocol.Response{}
	trans.hdlr.Handle(resp, req)
	return resp, nil
}

func (trans *memoryTransport) Serve() error {
	go trans.loop()
	return nil
}

func (trans *memoryTransport) loop() {
	for req := range trans.reqs {
		resp := &protocol.Response{}
		trans.hdlr.Handle(resp, req)
	}
}

func (trans *memoryTransport) AddHandler(proto protocol.Protocol, hdlr handler.Handler) {
	trans.proto = proto
	trans.hdlr = hdlr
}

func (t *memoryTransport) Listen() error {
	return nil
}
