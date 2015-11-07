package transport

type InMemoryTransport struct {
	reqs  chan Request
	resps [][]byte
}

func NewInMemoryTransport(reqs [][]byte, resps [][]byte) *InMemoryTransport {
	ch := make(chan Request, len(reqs)+1)
	for _, b := range reqs {
		ch <- Request{b, &DummyResponseWriter{}}
	}
	return &InMemoryTransport{
		ch,
		resps,
	}
}

func (trans *InMemoryTransport) Receive() <-chan Request {
	return trans.reqs
}

func (trans *InMemoryTransport) Send(endpoint string, data []byte) ([]byte, error) {
	return trans.resps[0], nil
}

func (t *InMemoryTransport) Listen(e string) {}

type DummyResponseWriter struct {
	Data  interface{}
	Error error
}

func (rw *DummyResponseWriter) Write(data interface{}) error {
	rw.Data = data
	return nil
}

func (rw *DummyResponseWriter) WriteError(err error) error {
	rw.Error = err
	return nil
}
