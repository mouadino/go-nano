package main

import (
	"fmt"

	"github.com/mouadino/go-nano/interfaces"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/reflection"
	"github.com/mouadino/go-nano/transport"
)

// TODO: Dependency Injection for transporter, serializer ...
// TODO: Service interface.
func Main(service interface{}) {
	app := NewApplication(
		service,
		protocol.NewJSONRPCProtocol(transport.NewHTTPTransport("127.0.0.1:8080")),
	)
	app.Serve()
}

type Application struct {
	protocol interfaces.Protocol
	handler  interfaces.Handler
	svc      interface{}
}

func NewApplication(svc interface{}, p interfaces.Protocol) *Application {
	return &Application{
		svc:      svc,
		protocol: p,
		handler:  reflection.FromStruct(svc),
	}
}

func (app *Application) Serve() {
	// TODO: goroutine Pool.
	for {
		resp, req := app.protocol.ReceiveRequest()
		fmt.Printf("%s -> %s\n", req, resp)
		// TODO: In a go routine.
		app.handler.Handle(resp, req)
	}
}

// TODO: Remove me !
func main() {}
