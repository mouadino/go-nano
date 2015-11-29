/*

Usage:

		type Upper struct {}

		func (Upper) Handle(rw protocol.ResponseWriter, req protocol.Request) {
			rw.Set(strings.ToUpper(text))
		}

		server := DefaultServer(proto)
		server.Register("Upper", Upper{}, zkAnnouncer)

		server.Serve()

*/
package nano

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/handler/middleware"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
)

// Hook represent a callback to call.
type Hook func() error

type hooks []Hook

func (hs hooks) Call() error {
	for _, h := range hs {
		err := h()
		if err != nil {
			return err
		}
	}
	return nil
}

// Server represent an RPC server.
type Server struct {
	trans   transport.Transport
	proto   protocol.Protocol
	hdlr    handler.Handler
	onStart hooks
	onStop  hooks
}

// DefaultServer returns a new nano.Server using default configuration.
// Default configuration include HTTP transport, JSONRPC as protocol and different
// middlewares including logger, tracing and recovery.
func DefaultServer(svc interface{}) *Server {
	trans := transport.NewHTTPTransport()
	server := CustomServer(
		handler.Reflect(svc),
		trans,
		jsonrpc.NewJSONRPCProtocol(trans, serializer.JSONSerializer{}),
		middleware.NewRecoverMiddleware(log.New(), true, 8*1024),
		middleware.NewTraceMiddleware(),
		middleware.NewLoggerMiddleware(log.New()),
	)
	if svc, ok := svc.(handler.Startable); ok {
		server.OnStart(Hook(svc.NanoStart))
		server.OnStop(Hook(svc.NanoStop))
	}
	return server
}

// CustomServer returns a new nano.Server customized with specific transport, protocol and middelwares.
func CustomServer(hdlr handler.Handler, trans transport.Transport, proto protocol.Protocol, middlewares ...handler.Middleware) *Server {
	hdlr = middleware.Chain(
		hdlr,
		middlewares...,
	)
	return &Server{
		trans: trans,
		proto: proto,
		hdlr:  hdlr,
	}
}

// OnStart add a hook to be called when starting the server.
func (s *Server) OnStart(h ...Hook) {
	s.onStart = append(s.onStart, h...)
}

// OnStop add a hook to be called when stoping the server.
func (s *Server) OnStop(h ...Hook) {
	s.onStop = append(s.onStop, h...)
}

// ListenAndServe listens on transport addr (if there is any) and then
// start handling requests from transport.
func (s *Server) ListenAndServe() error {
	s.trans.Listen()
	err := s.onStart.Call()
	if err != nil {
		return err
	}
	go s.loop()
	s.waitForTermination()
	return s.onStop.Call()
}

func (s *Server) loop() {
	for {
		resp, req, err := s.proto.Receive()
		if err != nil {
			log.Errorf("receive failed: %s", err)
			continue
		}
		go s.hdlr.Handle(resp, req)
	}
}

func (s *Server) waitForTermination() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Print("Received SIGTERM, exiting ...")
	}
}

// Announce register the RPC server in the passed announcer system under the given name with
// specific metadata.
func (s *Server) Announce(name string, serviceMeta discovery.ServiceMetadata, announcer discovery.Announcer) {
	s.OnStart(Hook(func() error {
		trans := s.trans.(transport.Listener)
		instance := discovery.NewInstance(
			discovery.NewServiceMetadata(
				trans.Addr(),
				serviceMeta,
			),
		)
		return announcer.Announce(name, instance)
	}))
}
