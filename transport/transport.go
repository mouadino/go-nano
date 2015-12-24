/*
Package transport defines logic to send and receive RPC requests and responses.
*/
package transport

import "io"

// Transport interface define behaviour of transport logic.
type Transport interface {
	Listen() error
	Receive() <-chan Request
	Send(string, io.Reader) ([]byte, error)
}

// Addresser interface should be implemented by services that support addressing.
// This interface is usually used with discovery to announce an RPC server with given transport.
type Addresser interface {
	Addr() string
}

// ResponseWriter represents interface that transport use to write response to the wire.
type ResponseWriter interface {
	Write(interface{}) error
}

// Request represents a incomming RPC request.
type Request struct {
	Body interface{}
	Resp ResponseWriter
}
