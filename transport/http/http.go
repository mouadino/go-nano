package http

import (
	"fmt"
	"net/http"

	"github.com/mouadino/go-nano/protocol"
)

const rpcPath = "/rpc/"

// NewSender returns a new RPC client that can send RPC requests using
// HTTP as a transport.
func NewSender(proto protocol.Protocol, opts ...func(*Sender)) *Sender {
	s := &Sender{
		proto:  proto,
		client: &http.Client{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func getContentType(proto protocol.Protocol) string {
	return fmt.Sprintf("application/%s", proto.String())
}

func createURL(svcName string) string {
	return rpcPath + svcName
}
