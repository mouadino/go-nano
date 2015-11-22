package transport

// TODO: Start/Stop
// TODO: Client side listen ?
type Transport interface {
	Listen() error
	Receive() <-chan Request
	Send(string, []byte) ([]byte, error)
}

// TODO: Addresser ?
type Listener interface {
	Addr() string
}

type ResponseWriter interface {
	Write(interface{}) error
	// TODO: Write should be only called once !?
}

type Request struct {
	Body interface{}
	Resp ResponseWriter
}
