package transport

type InMemoryTransport struct {
	ch chan []byte
}

func NewInMemoryTransport() *InMemoryTransport {
	// TODO: 10 look like the number of concurrent request that we can do !?
	return &InMemoryTransport{
		ch: make(chan []byte, 10),
	}
}

func (t *InMemoryTransport) Receive() <-chan []byte {
	return t.ch
}

func (t *InMemoryTransport) Send(endpoint string, data []byte) {
	t.ch <- data
}
