package nano

import (
	"fmt"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/reflection"
	"github.com/mouadino/go-nano/transport"
)

// TODO: Dependency Injection for transporter, serializer ...
// TODO: Service interface.
func Main(service interface{}) {
	t := transport.NewHTTPTransport()
	app := NewApplication(
		service,
		t,
		protocol.NewJSONRPCProtocol(t),
	)
	app.Serve()
}

type Application struct {
	transport transport.Transport
	protocol  protocol.Protocol
	handler   handler.Handler
	svc       interface{}
}

func NewApplication(svc interface{}, t transport.Transport, p protocol.Protocol) *Application {
	return &Application{
		svc:       svc,
		transport: t,
		protocol:  p,
		handler:   reflection.FromStruct(svc),
	}
}

func (app *Application) Serve() {
	// TODO: goroutine Pool.
	go app.transport.Listen("127.0.0.1:8080")
	for {
		resp, req := app.protocol.ReceiveRequest()
		fmt.Printf("%s -> %s\n", req, resp)
		go app.handler.Handle(resp, req)
	}
}
