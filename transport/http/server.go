package http

import (
	"net"
	"net/http"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
)

// Server represents a transport.Server using HTTP as transport.
type Server struct {
	httpServer  *http.Server
	mux         *http.ServeMux
	httpHandler http.Handler
}

// NewServer returns a new RPC Server that can serve RPC requests using
// HTTP as a transport.
func NewServer(hdlr handler.Handler) *Server {
	return &Server{
		mux:         http.NewServeMux(),
		httpHandler: newRPCHandler(hdlr),
	}
}

func (srv *Server) AddProtocol(proto protocol.Protocol) error {
	return srv.httpHandler.protos.Add(proto)
}

func (srv *Server) ListenAndServe(ln net.Listener) error {
	srv.mux.Handle(rpcPath, srv.httpHandler)
	return srv.httpServer.Serve(ln)
}
