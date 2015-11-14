package transport

type Transport interface {
	Listen(string)
	Receive() <-chan Request
	Send(string, []byte) ([]byte, error)
}

type Listener interface {
	Addr() string
}

type ResponseWriter interface {
	Write(interface{}) error
}

type Request struct {
	Body []byte
	Resp ResponseWriter
}
