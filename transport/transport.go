package transport

type Transport interface {
	Listen() error
	Receive() <-chan Request
	Send(string, []byte) ([]byte, error)
}

type Listener interface {
	Addr() string
}

type ResponseWriter interface {
	Write([]byte) error
	// TODO: Write should be only called once !?
}

type Request struct {
	Body []byte
	Resp ResponseWriter
}
