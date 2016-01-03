/*
Package nano is a toolkit for building services in golang. It's build on exiting mature golang
library and best practices to give it's user a nice (I am biased) api on how to create and manage
services in golang.

A usual service written with Nano is composed of a "handler" which is where business logic reside,
on top of that developers can use/write a "middlewares" which wraps a "handler" to implement application logic
e.g. logging, tracing, auth, rate limiting ... . that's from producer side, from consumer side Nano contains
a client to talk to services with client extensions which are same as middlewares but for client side, example
of such extensions is retry, backoff, circuit breaker ... .

A "handler" can be exposed using multiple transport and protocol, this make services written in Nano usable from other
services written using different stacks.

To give you more concrete description on how Nano can be used, we will start by writting a simple echo service.

Nano was designed in a way that hide the complexity of creating a service by letting it's users focus more on the
business logic of your service which is usually where you should start when writing a service.

In Nano word, business logic live inside what we call a "handler". To create a handler you can either create a simple
golang struct and let Nano reflect it:

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

Now that we have our business logic done, the next step it to expose it as a service. With Nano you can expose your
handler using different transport/protocol (multiple at the same time if necessary), for our simple example we will expose
our handler using HTTP transport and JSON/RPC protocol:

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

Something that is also very common when exposing a service is to be able to register it under a discovery system (like zookeeper),
in this case our main function should be:

		func main() {
			serv := server.New(jsonrpc.New(http.New()))
			serv.Register("echo", hdlr)
			serv.ServeAndAnnounce(zookeeper.New([]string{"127.0.0.1:2181"}))
		}

This make our handler available under the name "echo".

Now to consume our service using Nano, we will create a client and call our echo function:

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
			// Or c.Call("echo.echo", ...) If reflection was used to create the handler.
			if err != nil {
				panic(err)
			}
			fmt.Printf("Echo returned %s", msg)
		}

If service discovery was used there will be no need to have a hardcoded echo service endpoint, instead
our client setup will look like this:

		func main() {
			lb := loadbalancer.New(
				zookeeper.New([]string{"127.0.0.1:2181"}),
				loadbalancer.NewRoundRobin(),
			)

			c := client.New(echo, jsonrpc.New(http.New()), lb)

			// Then use client like above.

		}

*/
package nano
