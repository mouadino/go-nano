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

type Hook func() error

type Hooks []Hook

func (hs Hooks) Call() error {
	for _, h := range hs {
		err := h()
		if err != nil {
			return err
		}
	}
	return nil
}

type Server struct {
	trans   transport.Transport
	proto   protocol.Protocol
	hdlr    handler.Handler
	onStart Hooks
	onStop  Hooks
}

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

func (s *Server) OnStart(h ...Hook) {
	s.onStart = append(s.onStart, h...)
}

func (s *Server) OnStop(h ...Hook) {
	s.onStop = append(s.onStop, h...)
}

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
		resp, req := s.proto.ReceiveRequest()
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

func (s *Server) Announce(name string, serviceMeta discovery.ServiceMetadata, announcer discovery.Announcer) {
	s.OnStart(Hook(func() error {
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
	}))
}
