package transport

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type HTTPResponseWriter struct {
	sent chan struct{}
	resp http.ResponseWriter
}

func (w *HTTPResponseWriter) Write(data interface{}) error {
	_, err := w.resp.Write(data.([]byte))
	if err != nil {
		return err
	}
	w.sent <- struct{}{}
	return nil
}

type HTTPTransport struct {
	mux  *http.ServeMux
	reqs chan Request
}

func NewHTTPTransport() *HTTPTransport {
	return &HTTPTransport{
		mux:  http.NewServeMux(),
		reqs: make(chan Request),
	}
}

func (trans *HTTPTransport) Listen(address string) {
	trans.mux.HandleFunc("/rpc/", trans.handler)
	listner, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Listening failed: ", err)
	}
	log.Info("Listening on ", listner.Addr())
	http.Serve(listner, trans.mux)
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
	// TODO: Timeout !?
	<-resp.sent
}

func (trans *HTTPTransport) Send(endpoint string, body []byte) ([]byte, error) {
	endpoint = fmt.Sprintf("%s/rpc/", endpoint)
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
