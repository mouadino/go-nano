package memory

import (
	"io"

	"github.com/mouadino/go-nano/transport"
)

type memoryTransport struct {
	reqs  chan transport.Request
	resps [][]byte
}

func New(reqs [][]byte, resps [][]byte) transport.Transport {
	ch := make(chan transport.Request, len(reqs)+1)
	for _, b := range reqs {
		ch <- transport.Request{b, &DumpResponseWriter{}}
	}
	return &memoryTransport{
		ch,
		resps,
	}
}

func (trans *memoryTransport) Receive() <-chan transport.Request {
	return trans.reqs
}

func (trans *memoryTransport) Send(endpoint string, r io.Reader) ([]byte, error) {
	return trans.resps[0], nil
}

func (t *memoryTransport) Listen() error {
	return nil
}

func (t *memoryTransport) Endpoint() string {
	return "<memory>"
}

type DumpResponseWriter struct {
	Data interface{}
}

func (rw *DumpResponseWriter) Write(data interface{}) error {
	rw.Data = data
	return nil
}
