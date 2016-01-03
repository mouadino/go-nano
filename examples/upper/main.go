package main

import (
	"flag"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/mouadino/go-nano/discovery/zookeeper"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/handler/middleware"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/server"
	"github.com/mouadino/go-nano/transport/amqp"
	"github.com/mouadino/go-nano/transport/http"
)

var zkHost = flag.String("zookeeper", "127.0.0.1:2181", "Zookeeper location")
var rmqHost = flag.String("rabbitmq", "amqp://127.0.0.1:5672", "RabbitMQ location")
var logger = log.New()
var ms = []handler.Middleware{
	middleware.NewLogger(logger),
	middleware.NewRecover(logger, true, 8*1024),
	middleware.NewTrace(),
}

type upperService struct{}

func (upperService) Upper(s string) (string, error) {
	return strings.ToUpper(s), nil
}

func main() {
	flag.Parse()

	zkAnnouncer := zookeeper.New(
		[]string{*zkHost},
	)

	httpServ := server.New(jsonrpc.New(http.New()))
	httpServ.Register("upper", upperService{}, ms...)

	go func() {
		if err := httpServ.ServeAndAnnounce(zkAnnouncer); err != nil {
			panic(err)
		}
	}()

	amqpServ := server.New(jsonrpc.New(amqp.New(*rmqHost)))
	amqpServ.Register("upper", upperService{}, ms...)

	amqpServ.Serve()
}
