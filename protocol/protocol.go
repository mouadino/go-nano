package protocol

import "github.com/mouadino/go-nano/header"

type Protocol interface {
	// TODO: SendRequest(string, *Request) (Response, error)
	SendRequest(string, *Request) (interface{}, error)
	// TODO: How to return an error ?
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

	// TODO: s/Write/Set ?
	Write(interface{}) error
	WriteError(err error) error
}
