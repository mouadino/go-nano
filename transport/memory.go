package transport

type InMemoryTransport struct {
	reqs  chan Request
	resps [][]byte
}

func NewInMemoryTransport(reqs [][]byte, resps [][]byte) *InMemoryTransport {
	ch := make(chan Request, len(reqs)+1)
	for _, b := range reqs {
		ch <- Request{b, &DumpResponseWriter{}}
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

func (t *InMemoryTransport) Listen() error {
	return nil
}

func (t *InMemoryTransport) Endpoint() string {
	return "<memory>"
}

type DumpResponseWriter struct {
	Data interface{}
}

func (rw *DumpResponseWriter) Write(data []byte) error {
	rw.Data = data
	return nil
}
