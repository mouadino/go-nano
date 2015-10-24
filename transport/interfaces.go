package transport

import "github.com/mouadino/go-nano/header"

type Transport interface {
	Listen(string)
	Receive() <-chan Data
	Send(endpoint string, data []byte) (ResponseReader, error)
}

type ResponseWriter interface {
	Header() header.Header

	Write(interface{}) error
}

// TODO: Rename me
type Data struct {
	Body []byte
	Resp ResponseWriter
}

type ResponseReader interface {
	Read() ([]byte, error)
}
