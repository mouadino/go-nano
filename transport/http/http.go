package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/transport"
)

const RPCPath = "/rpc/"

type HTTPTransport struct {
	mux       *http.ServeMux
	reqs      chan transport.Request
	addr      string
	listening bool
}

// New creates a new HTTP transport.
func New() transport.Transport {
	return &HTTPTransport{
		mux:  http.NewServeMux(),
		reqs: make(chan transport.Request),
	}
}

func (trans *HTTPTransport) Listen() error {
	trans.mux.HandleFunc(RPCPath, trans.handler)
	// TODO: Get external IP.
	listner, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	trans.addr = fmt.Sprintf("http://%s", listner.Addr().String())
	log.Info("Listening on ", trans.addr)
	go http.Serve(listner, trans.mux)
	trans.listening = true
	return nil
}

func (trans *HTTPTransport) Addr() string {
	return trans.addr
}

func (trans *HTTPTransport) handler(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Errorf("Transport error: %s", err)
		return
	}
	resp := HTTPResponseWriter{
		make(chan struct{}),
		rw,
	}
	trans.reqs <- transport.Request{
		Body: body,
		Resp: &resp,
	}
	// Wait here else http.ResponseWriter became invalid.
	// TODO: Timeout !?, http.Hijacker ?
	<-resp.sent
}

func (trans *HTTPTransport) Send(endpoint string, body io.Reader) ([]byte, error) {
	endpoint += RPCPath
	// TODO: content-type doesn't belong here.
	resp, err := http.Post(endpoint, "application/json-rpc", body)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (trans *HTTPTransport) Receive() <-chan transport.Request {
	return trans.reqs
}
