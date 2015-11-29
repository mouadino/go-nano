package protocol

import "github.com/mouadino/go-nano/header"

type Sender interface {
	Send(string, *Request) (*Response, error)
}

type Receiver interface {
	Receive() (ResponseWriter, *Request, error)
}

type Protocol interface {
	Sender
	Receiver
}

type Params map[string]interface{}

type Request struct {
	Method string
	Params Params
	Header header.Header
}

type Response struct {
	Body   interface{}
	Error  error
	Header header.Header
}

// TODO: s/ResponseWriter/Responder/ ?
type ResponseWriter interface {
	Header() header.Header

	Set(interface{}) error
	SetError(err error) error
}
