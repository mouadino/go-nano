package client

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mouadino/go-nano/client/extension"
	"github.com/mouadino/go-nano/transport"
)

/*
URI:

- http+jsonrpc://127.0.0.1:8089/rpc/demo
- tcp+lymph://127.0.0.0.1:309180/
- amqp+jsonrpc//127.0.0.1:5432/demo
- zk://127.0.0.1:2181/demo
*/
type uri struct {
	trans    string
	proto    string
	endpoint string
}

func NewGeneric(endpoint string, exts ...extension.Extension) (*Client, error) {
	// TODO: Need to call zookeeper first
	sender, err := getSender(endpoint)
	if err != nil {
		return nil, err
	}
	client := New(endpoint, sender, exts...)
	return &client, nil
}

func getSender(endpoint string) (transport.Sender, error) {
	// TODO: Need to call discovery in case it's needed.
	_, err := parseURI(endpoint)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Parse an URI in the form <trans>+<proto>://<location>
// Example:
//
//   http+jsonrpc://127.0.0.1:8080/
//
func parseURI(endpoint string) (*uri, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	scheme := strings.SplitN(u.Scheme, "+", 2)
	if len(scheme) != 2 {
		return nil, fmt.Errorf("malformed uri scheme: %s", endpoint)
	}

	uri := &uri{
		trans:    scheme[0],
		proto:    scheme[1],
		endpoint: endpoint,
	}

	return uri, nil
}
