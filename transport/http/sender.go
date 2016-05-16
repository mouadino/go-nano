package http

import (
	"bytes"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"

	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
)

// HTTPTimeout configure timeout in the underline HTTP client.
func HTTPTimeout(timeout time.Duration) func(*Sender) {
	return func(s *Sender) {
		s.client.Timeout = timeout
	}
}

// Sender represents a transport.Sender using HTTP as transport.
type Sender struct {
	proto  protocol.Protocol
	client *http.Client
}

// Send a raw request to given endpoint using HTTP as transport.
func (snd *Sender) Send(endpoint string, ctx context.Context, protoReq *protocol.Request) (*protocol.Response, error) {
	url := createURL(endpoint)
	req, err := snd.makeRequest(url, protoReq)
	if err != nil {
		return nil, err
	}

	resp, err := ctxhttp.Do(ctx, snd.client, req)
	if err != nil {
		return nil, err
	}
	return snd.proto.DecodeResponse(resp.Body, header.Header(resp.Header))
}

func (snd *Sender) makeRequest(url string, protoReq *protocol.Request) (*http.Request, error) {
	body, err := snd.proto.EncodeRequest(protoReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	snd.setHeaders(req, protoReq)
	return req, nil
}

func (snd *Sender) setHeaders(req *http.Request, protoReq *protocol.Request) {
	for k, vs := range protoReq.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}

	contentType := getContentType(snd.proto)
	req.Header.Set("Content-Type", contentType)
}
