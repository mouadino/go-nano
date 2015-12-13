package main

import (
	"flag"
	"strings"

	"github.com/mouadino/go-nano/discovery/zookeeper"
	"github.com/mouadino/go-nano/handler/middleware"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/server"
	"github.com/mouadino/go-nano/transport/amqp"
	"github.com/mouadino/go-nano/transport/http"
)

var zkHost = flag.String("zookeeper", "127.0.0.1:2181", "Zookeeper location")
var rmqHost = flag.String("rabbitmq", "amqp://127.0.0.1:5672", "RabbitMQ location")

type upperService struct{}

func (upperService) Upper(s string) string {
	return strings.ToUpper(s)
}

func main() {
	flag.Parse()

	zkAnnouncer := zookeeper.New(
		[]string{*zkHost},
	)

	serv := server.New(jsonrpc.New(http.New()))
	serv.Register("upper", upperService{}, middleware.Defaults...)

	if err := serv.ServeAndAnnounce(zkAnnouncer); err != nil {
		panic(err)
	}

	serv = server.New(jsonrpc.New(amqp.New(*rmqHost)))
	serv.Register("upper", upperService{}, middleware.Defaults...)

	serv.Serve()

	server.Wait()
}
