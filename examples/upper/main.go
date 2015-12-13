package main

import (
	"strings"

	"github.com/mouadino/go-nano/discovery/zookeeper"
	"github.com/mouadino/go-nano/handler/middleware"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/server"
	"github.com/mouadino/go-nano/transport/amqp"
	"github.com/mouadino/go-nano/transport/http"
)

type upperService struct{}

func (upperService) Upper(s string) string {
	return strings.ToUpper(s)
}

func main() {
	zkAnnouncer := zookeeper.New(
		[]string{"127.0.0.1:2181"},
	)

	serv := server.New(jsonrpc.New(http.New()))
	serv.Register("upper", upperService{}, middleware.Defaults...)

	if err := serv.ServeAndAnnounce(zkAnnouncer); err != nil {
		panic(err)
	}

	serv = server.New(jsonrpc.New(amqp.New("amqp://127.0.0.1:5672")))
	serv.Register("upper", upperService{}, middleware.Defaults...)

	serv.Serve()

	server.Wait()
}
