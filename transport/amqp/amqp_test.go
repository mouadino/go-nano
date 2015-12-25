// +build integration
package amqp

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/mouadino/go-nano/transport"
)

const (
	testExchangeName = "nano-test"
	amqpURI          = "amqp://guest:guest@localhost:5672/"
)

func chanRead(ch <-chan transport.Request, timeout time.Duration) (transport.Request, error) {
	select {
	case res := <-ch:
		return res, nil
	case <-time.After(2 * time.Second):
		return transport.Request{}, errors.New("Timeout")
	}
}

func TestAMQPTransport(t *testing.T) {
	firstEndpoint := "first-test-client"
	firstTrans := New(amqpURI, Exchange(testExchangeName), QueueName(firstEndpoint))
	firstTrans.Listen()

	secondEndpoint := "second-test-client"
	secondTrans := New(amqpURI, Exchange(testExchangeName), QueueName(secondEndpoint))
	secondTrans.Listen()

	body := "Hello World!"
	type RespData struct {
		body []byte
		err  error
	}
	respCh := make(chan RespData, 1)
	go func() {
		body, err := firstTrans.Send(secondEndpoint, strings.NewReader(body))
		respCh <- RespData{body, err}
	}()

	req, err := chanRead(secondTrans.Receive(), 2*time.Second)
	if err != nil {
		t.Fatalf("Didn't receive any request after 2 second")
	}

	b, ok := req.Body.([]byte)
	if !ok {
		t.Fatalf("request body is not []byte")
	}
	if string(b) != body {
		t.Errorf("request body doesn't match want %v, got %v", body, req.Body)
	}

	req.Resp.Write([]byte(body))

	resp := <-respCh
	if resp.err != nil {
		t.Fatalf("%s", resp.err)
	}

	if string(resp.body) != body {
		t.Errorf("response body didn't match, want %s got %s", body, resp.body)
	}

}
