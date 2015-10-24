package protocol

import "github.com/mouadino/go-nano/transport"

type Protocol interface {
	SendRequest(string, *Request) (interface{}, error)
	ReceiveRequest() (transport.ResponseWriter, *Request)
	// TODO: SendError !?
}

type Params map[string]interface{}

// TODO: Is this specific to JSON-RPC !?
type Request struct {
	Method string
	Params Params
	// TODO: Headers header
}
