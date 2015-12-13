package transport

import "io"

type Transport interface {
	Listen() error
	Receive() <-chan Request
	Send(string, io.Reader) ([]byte, error)
}

type Addresser interface {
	Addr() string
}

type ResponseWriter interface {
	Write(interface{}) error
}

type Request struct {
	Body interface{}
	Resp ResponseWriter
}
