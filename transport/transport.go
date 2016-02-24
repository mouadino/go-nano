/*
Package transport defines logic to send and receive RPC requests and responses.
*/
package transport

import (
	"net"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	"golang.org/x/net/context"
)

type Sender interface {
	Send(string, context.Context, *protocol.Request) (*protocol.Response, error)
}

// Server interface represents an RPC server.
type Server interface {
	AddHandler(proto protocol.Protocol, hdlr handler.Handler)
	Serve() error
}

// Listener interface should be implemented by services that support addressing.
// This interface is usually used with discovery to announce an RPC server with
// given server.
type Listener interface {
	Listen(net.Listener) error
	ListenAndServe(net.Listener) error
}

// Startable interface should be implemented by transport that are need to
// setup (or clean up) some resources before beign used.
type Startable interface {
	Start() error
	Stop() error
}
