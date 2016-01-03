/*
Package nano is a toolkit for building services in golang. It's build on exiting mature golang
library and best practices to give it's user a familiar setup on how to create and manage services
in golang.

A usual service written with "nano" is composed of a "handler" which is where business logic reside,
on top of that we can define "middlewares" which wraps a "handler" to implement application logic
dynamically e.g. logging, tracing, auth ... . A "handler" can be exposed using different transport
and protocol, this allow a service written in "nano" to talk to other services written using different
stacks.

Here we will define a high overview on how "nano" services are structured, to do that we will be writing
a simple echo service..

you first start by defining a "handler" which can be done either by using reflection:

		package main

		import (
			"github.com/mouadino/go-nano/handler"
		)

		type Echo struct {}

		func (Echo) Echo(msg string) (string, error) {
			return msg, nil
		}

		var hdlr = handler.New(Echo{})

Or by implementing the handler.Handler interface directly:

		package main

		import (
			"github.com/mouadino/go-nano/handler"
			"github.com/mouadino/go-nano/protocol"
		)

		var hdlr = handler.HandlerFunc(func(rw protocol.ResponseWriter, req *protocol.Request) {
			msg := req.Params["_0"].(string)
			rw.Set(msg)
		}

Next we will expose our handler using HTTP transport and JSON/RPC protocol:

		package main

		import (
			"github.com/mouadino/go-nano/server"
			"github.com/mouadino/go-nano/protocol/jsonrpc"
			"github.com/mouadino/go-nano/transport/http"
		)

		func main() {
			serv := server.New(jsonrpc.New(http.New()))

			serv.Register("echo", hdlr)

			serv.Serve()
		}

To use service discovery you can replace main function with:

		func main() {
			serv := server.New(jsonrpc.New(http.New()))

			serv.Register("echo", hdlr)

			zk := zookeeper.New([]string{"127.0.0.1:2181"})
			serv.ServeAndAnnounce(zk)
		}

This make our hanlder available under the name "echo".

Now from client side we can talk to our "echo" service:

		package main

		import (
			"fmt"

			"github.com/mouadino/go-nano/client"
			"github.com/mouadino/go-nano/protocol/jsonrpc"
			"github.com/mouadino/go-nano/transport/http"
		)

		func main() {
			c := client.New("http://127.0.0.1:23765", jsonrpc.New(http.New()))

			msg, err := c.Call("echo", "Hello, World!")

			if err != nil {
				panic(err)
			}

			fmt.Printf("Echo returned %s", msg)
		}

If service discovery was used we will not have to hardcode the echo service endpoint, instead
our client setup will look like this:

		func main() {
			zk := zookeeper.New([]string{"127.0.0.1:2181"})
			lb := loadbalancer.New(zk, loadbalancer.NewRoundRobin())

			c := client.New(echo, jsonrpc.New(http.New()), lb)

			// Then use client like above.

		}

*/
package nano
