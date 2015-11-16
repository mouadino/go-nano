package transport

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const RPCPath = "/rpc/"

type HTTPResponseWriter struct {
	sent chan struct{}
	resp http.ResponseWriter
}

func (w *HTTPResponseWriter) Write(data []byte) error {
	_, err := w.resp.Write(data)
	if err != nil {
		return err
	}
	w.sent <- struct{}{}
	return nil
}

type HTTPTransport struct {
	mux  *http.ServeMux
	reqs chan Request
	addr string
}

func NewHTTPTransport() Transport {
	return &HTTPTransport{
		mux:  http.NewServeMux(),
		reqs: make(chan Request),
	}
}

func (trans *HTTPTransport) Listen() error {
	trans.mux.HandleFunc(RPCPath, trans.handler)
	listner, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	trans.addr = fmt.Sprintf("http://%s", listner.Addr().String())
	log.Info("Listening on ", trans.addr)
	go http.Serve(listner, trans.mux)
	return nil
}

func (trans *HTTPTransport) Addr() string {
	return trans.addr
}

func (trans *HTTPTransport) handler(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Error("Transport error: %s", err)
		return
	}
	resp := HTTPResponseWriter{
		make(chan struct{}),
		rw,
	}
	trans.reqs <- Request{
		Body: body,
		Resp: &resp,
	}
	// Wait here else http.ResponseWriter became invalid.
	// TODO: Timeout !?, http.Hijacker ?
	<-resp.sent
}

func (trans *HTTPTransport) Send(endpoint string, body []byte) ([]byte, error) {
	endpoint += RPCPath
	// TODO: content-type doesn't belong here.
	resp, err := http.Post(endpoint, "application/json-rpc", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (trans *HTTPTransport) Receive() <-chan Request {
	return trans.reqs
}
