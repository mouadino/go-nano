/*
package server define how to serve different handlers with different
namespaces.

Usage:

		type Upper struct {}

		func (Upper) Handle(rw protocol.ResponseWriter, req *protocol.Request) {
			text := req.Params["text"].(string)
			rw.Set(strings.ToUpper(text))
		}

		serv := server.New(jsonrpc.New(http.New())).build()
		serv.Register("Upper", Upper{})

		_ = serv.Serve()

*/
package server

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/handler/middleware"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

// Server represents an RPC server.
type Server struct {
	proto protocol.Protocol
	mux   handlersMux
	metas map[string]map[string]interface{}
}

func New(proto protocol.Protocol) *Server {
	return &Server{
		proto: proto,
		mux:   handlersMux{make(map[string]handler.Handler)},
		metas: make(map[string]map[string]interface{}),
	}
}

// Register given handler under name.
func (s *Server) Register(name string, svc interface{}, ms ...handler.Middleware) error {
	hdlr := middleware.Chain(handler.Reflect(svc), ms...)
	return s.mux.Register(name, hdlr)
}

// Register given handler under name with given metadata.
func (s *Server) RegisterWithMetadata(name string, svc interface{}, meta map[string]interface{}, ms ...handler.Middleware) error {
	err := s.Register(name, svc, ms...)
	if err != nil {
		return err
	}
	s.metas[name] = meta
	return nil
}

// Serve listens on transport addr (if there is any) and then
// start handling requests from transport.
func (s *Server) Serve() {
	trans := s.proto.Transport()
	trans.Listen()
	go s.loop()
}

func (s *Server) ServeAndAnnounce(an discovery.Announcer) error {
	s.Serve()
	return s.announce(an)
}

func (s *Server) announce(an discovery.Announcer) error {
	addr, ok := s.proto.Transport().(transport.Addresser)
	if !ok {
		return errors.New("can only announce transport of type transport.Addresser")
	}

	for _, name := range s.mux.Names() {
		meta := discovery.NewServiceMetadata(addr.Addr(), s.metas[name])
		instance := discovery.NewInstance(meta)
		err := an.Announce(name, instance)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) loop() {
	for {
		resp, req, err := s.proto.Receive()
		if err != nil {
			log.Errorf("transport receive failed: %s", err)
			continue
		}
		if err != nil {
			log.Errorf("code failed to decode: %s", err)
			continue
		}
		go s.mux.Handle(resp, req)
	}
}

func Wait() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Print("Received SIGTERM, exiting ...")
	}
}
