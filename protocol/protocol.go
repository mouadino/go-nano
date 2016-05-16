package protocol

import (
	"io"

	"github.com/mouadino/go-nano/header"
)

type Protocol interface {
	EncodeRequest(*Request) ([]byte, error)
	DecodeRequest(io.Reader, header.Header) (*Request, error)

	EncodeResponse(*Response) ([]byte, error)
	DecodeResponse(io.Reader, header.Header) (*Response, error)

	// String returns the name of the protocol, to be used in content-type and uri scheme.
	String() string
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
