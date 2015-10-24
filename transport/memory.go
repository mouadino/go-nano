package transport

import "github.com/mouadino/go-nano/header"

type InMemoryTransport struct {
	ch chan Data
}

func NewInMemoryTransport(data ...Data) *InMemoryTransport {
	ch := make(chan Data, len(data)+1)
	for _, d := range data {
		ch <- d
	}
	return &InMemoryTransport{ch}
}

func (t *InMemoryTransport) Receive() <-chan Data {
	return t.ch
}

func (t *InMemoryTransport) Send(endpoint string, data []byte) (ResponseReader, error) {
	return &DummyResponseReader{}, nil
}

func (t *InMemoryTransport) Listen(e string) {}

type DummyResponseReader struct{}

func (t *DummyResponseReader) Read() ([]byte, error) {
	return []byte{}, nil
}

type DummyResponseWriter struct{}

func (w *DummyResponseWriter) Write(data interface{}) error {
	return nil
}

func (w *DummyResponseWriter) Header() header.Header {
	return map[string][]string{} // FIXME: w.resp.Header()
}
