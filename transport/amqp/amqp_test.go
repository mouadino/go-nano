// +build integration
package amqp

import (
	"testing"

	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/jsonrpc"
)

const (
	testExchangeName = "nano-test"
	amqpURI          = "amqp://guest:guest@localhost:5672/"
)

type echoHandler struct{}

func (echoHandler) Handle(resp *protocol.Response, req *protocol.Request) {
	resp.Body = req.Params["_0"]
}

func TestAMQPTransport(t *testing.T) {
	firstEndpoint := "first-test-client"
	firstTrans := New(amqpURI, Exchange(testExchangeName), QueueName(firstEndpoint))
	firstTrans.AddHandler(jsonrpc.New(), echoHandler{})
	firstTrans.Listen()
	go firstTrans.Serve()

	secondEndpoint := "second-test-client"
	secondTrans := New(amqpURI, Exchange(testExchangeName), QueueName(secondEndpoint))
	secondTrans.AddHandler(jsonrpc.New(), echoHandler{})
	secondTrans.Listen()
	go secondTrans.Serve()

	req := &protocol.Request{
		Method: "",
		Params: protocol.Params{"_0": "foobar"},
	}

	resp, err := firstTrans.Send(secondEndpoint, req)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if resp.Body != "foobar" {
		t.Errorf("response body want %v, got %v", "foobar", resp.Body)
	}
}
