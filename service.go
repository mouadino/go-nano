package nano

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/reflection"
	"github.com/mouadino/go-nano/transport"
)

func Default(service interface{}) *Service {
	trans := transport.NewHTTPTransport()
	return Custom(
		service,
		trans,
		protocol.NewJSONRPCProtocol(trans),
	)
}

func Custom(svc interface{}, trans transport.Transport, proto protocol.Protocol) *Service {
	return &Service{
		svc:     svc,
		trans:   trans,
		proto:   proto,
		handler: reflection.FromStruct(svc),
	}
}

type Service struct {
	trans   transport.Transport
	proto   protocol.Protocol
	handler handler.Handler
	svc     interface{}
}

func (s *Service) ListenAndServe() {
	if s, ok := s.svc.(Startable); ok {
		err := s.NanoStart()
		if err != nil {
			log.Fatalf("Service failed to start: %s", err)
		}
		defer s.NanoStop()
	}
	// TODO: goroutine Pool.
	go s.trans.Listen(":0")
	go s.loop()
	s.waitForTermination()
}

func (s *Service) loop() {
	for {
		resp, req := s.proto.ReceiveRequest()
		log.Printf("%s -> %s\n", req, resp)
		go s.handler.Handle(resp, req)
	}
}

func (s *Service) waitForTermination() {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Print("Received SIGTERM, exiting ...")
	}
}
