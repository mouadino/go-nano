package protocol

import "github.com/mouadino/go-nano/header"

type Protocol interface {
	SendRequest(string, *Request) (interface{}, error)
	ReceiveRequest() (ResponseWriter, *Request)
}

type Params map[string]interface{}

type Request struct {
	Method string
	Params Params
	Header header.Header
}

type ResponseWriter interface {
	Header() header.Header

	Write(interface{}) error
	WriteError(err error) error
}
