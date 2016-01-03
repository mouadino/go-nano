package protocol

import (
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/transport"
)

type Sender interface {
	Send(string, *Request) (*Response, error)
}

type Receiver interface {
	Receive() (ResponseWriter, *Request, error)
}

type Protocol interface {
	Sender
	Receiver
	Transport() transport.Transport
}

type ProtocolV2 interface {
	EncodeRequest(Request) []byte
	DecodeRequest([]byte) Request

	EncodeResponse(Response) []byte
	DecodeResponse([]byte) Response

	EncodeError(error) []byte
	DecodeError([]byte) error

	// String returns the name of the protocol, to be used in content-type and uri scheme.
	String()
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
