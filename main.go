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
	t := transport.NewHTTPTransport()
	app := NewApplication(
		service,
		t,
		protocol.NewJSONRPCProtocol(t),
	)
	app.Serve()
}

type Application struct {
	transport interfaces.Transport
	protocol  interfaces.Protocol
	handler   interfaces.Handler
	svc       interface{}
}

func NewApplication(svc interface{}, t interfaces.Transport, p interfaces.Protocol) *Application {
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
		// TODO: In a go routine.
		go app.handler.Handle(resp, req)
	}
}

type remoteClient struct {
	endpoint string
	protocol interfaces.Protocol
}

func Client(endpoint string) *remoteClient {
	return &remoteClient{
		endpoint: endpoint,
		protocol: protocol.NewJSONRPCProtocol(transport.NewHTTPTransport()),
	}
}

func (c *remoteClient) Call(method string, params map[string]interface{}) (interface{}, error) {
	req := interfaces.Request{
		Method: method,
		Params: params,
	}
	resp, err := c.protocol.SendRequest(c.endpoint, &req)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TODO: Remove me !
func main() {}
