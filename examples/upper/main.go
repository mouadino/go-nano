package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"

	nano "github.com/mouadino/go-nano"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/handler/middleware"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
)

type echoService struct{}

func (echoService) NanoStart() error {
	log.Debug("Starting ...")
	return nil
}

func (echoService) NanoStop() error {
	log.Debug("Stopping ...")
	return nil
}

func (echoService) Upper(s string) string {
	return strings.ToUpper(s)
}

func main() {
	/*zkAnnouncer := discovery.DefaultZooKeeperAnnounceResolver(
		[]string{"127.0.0.1:2181"},
	)
	server := nano.DefaultServer(echoService{})
	server.Announce("upper", discovery.ServiceMetadata{}, zkAnnouncer)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}*/

	amqpTrans := transport.NewAMQPTransport("amqp://127.0.0.1:5672")
	server := nano.CustomServer(
		handler.Reflect(echoService{}),
		amqpTrans,
		jsonrpc.NewJSONRPCProtocol(amqpTrans, serializer.JSONSerializer{}),
		middleware.NewRecoverMiddleware(log.New(), true, 8*1024),
		middleware.NewTraceMiddleware(),
		middleware.NewLoggerMiddleware(log.New()),
	)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
