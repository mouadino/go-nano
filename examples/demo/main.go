package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/mouadino/go-nano/client"
	"github.com/mouadino/go-nano/client/extension"
	"github.com/mouadino/go-nano/discovery/loadbalancer"
	"github.com/mouadino/go-nano/discovery/zookeeper"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
	"github.com/mouadino/go-nano/transport/amqp"
	"github.com/mouadino/go-nano/transport/http"
)

var trans = flag.String("transport", "http", "transport to use e.g. http, amqp")
var zkHost = flag.String("zookeeper", "127.0.0.1:2181", "Zookeeper location")
var rmqHost = flag.String("rabbitmq", "amqp://127.0.0.1:5672", "RabbitMQ location")

func main() {
	flag.Parse()

	var cl client.Client
	if *trans == "http" {
		zk := zookeeper.New([]string{*zkHost})
		lb := loadbalancer.New(zk, loadbalancer.NewRoundRobin())

		// TODO: Add protocol.
		cl = client.New("upper", http.New(), lb, extension.NewTimeoutExt(1*time.Second))
	} else if *trans == "amqp" {
		cl = client.New("upper", amqp.New(*rmqHost), jsonrpc.New(), extension.NewTimeoutExt(2*time.Second))
		panic("amqp not supported")
	} else {
		panic("unknown transport")
	}

	c := time.Tick(1 * time.Second)
	i := 0
	for _ = range c {
		text := fmt.Sprintf("foo_%d", i)
		result, err := cl.Call("upper.upper", text)

		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("%s\n", result.(string))
		}
		i++
	}
}
