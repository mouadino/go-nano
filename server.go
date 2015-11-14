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
	"github.com/mouadino/go-nano/reflection"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
)

func DefaultServer(service interface{}) *Server {
	trans := transport.NewHTTPTransport()
	return CustomServer(
		service,
		trans,
		jsonrpc.NewJSONRPCProtocol(trans, serializer.JSONSerializer{}),
		middleware.NewRecoverMiddleware(log.New(), true, 8*1024),
		middleware.NewTraceMiddleware(),
		middleware.NewLoggerMiddleware(log.New()),
	)
}

func CustomServer(svc interface{}, trans transport.Transport, proto protocol.Protocol, middlewares ...handler.Middleware) *Server {
	handler := middleware.Chain(
		reflection.FromStruct(svc),
		middlewares...,
	)
	server := &Server{
		svc:     svc,
		trans:   trans,
		proto:   proto,
		handler: handler,
	}
	// FIXME: Not Good :(
	server.trans.Listen("127.0.0.1:0")
	return server
}

type Server struct {
	trans   transport.Transport
	proto   protocol.Protocol
	handler handler.Handler
	svc     interface{}
}

func (s *Server) ListenAndServe() {
	if svc, ok := s.svc.(Startable); ok {
		err := svc.NanoStart()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("Server failed to start")
		}
		defer svc.NanoStop()
	}

	go s.loop()
	s.waitForTermination()
}

func (s *Server) loop() {
	for {
		resp, req := s.proto.ReceiveRequest()
		go s.handler.Handle(resp, req)
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

func (s *Server) Announce(name string, serviceMeta discovery.ServiceMetadata, announcer discovery.Announcer) error {
	trans := s.trans.(transport.Listener)
	instance, err := discovery.NewInstance(
		discovery.NewServiceMetadata(
			trans.Addr(),
			serviceMeta,
		),
	)
	if err != nil {
		return err
	}
	return announcer.Announce(name, instance)
}
