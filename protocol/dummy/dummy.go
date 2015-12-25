package dummy

import (
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
	"github.com/mouadino/go-nano/transport/memory"
)

type ResponseWriter struct {
	Data         interface{}
	Error        error
	HeaderValues header.Header
}

func NewResponseRecorder() *ResponseWriter {
	return &ResponseWriter{}
}

func (rw *ResponseWriter) Header() header.Header {
	return rw.HeaderValues
}

func (rw *ResponseWriter) Set(data interface{}) error {
	rw.Data = data
	return nil
}

func (rw *ResponseWriter) SetError(err error) error {
	rw.Error = err
	return nil
}

type dummyProtocol struct {
	trans transport.Transport
	rw    *ResponseWriter
	req   *protocol.Request
}

func New(rw *ResponseWriter, req *protocol.Request) protocol.Protocol {
	return &dummyProtocol{
		trans: memory.New([][]byte{}, [][]byte{}),
		req:   req,
		rw:    rw,
	}
}

func (p *dummyProtocol) Receive() (protocol.ResponseWriter, *protocol.Request, error) {
	return p.rw, p.req, nil
}

func (p *dummyProtocol) Transport() transport.Transport {
	return p.trans
}

func (p *dummyProtocol) Send(e string, req *protocol.Request) (*protocol.Response, error) {
	// TODO: Dummy implementation
	return &protocol.Response{
		Body:  "",
		Error: nil,
	}, nil
}
