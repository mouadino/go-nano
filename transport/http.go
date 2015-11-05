package transport

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/mouadino/go-nano/header"
)

type ResponsePromise struct {
	ready chan struct{}
	resp  http.Response
	err   error
}

func (r *ResponsePromise) Read() ([]byte, error) {
	<-r.ready
	if r.err != nil {
		return []byte{}, r.err
	}
	defer r.resp.Body.Close()
	return ioutil.ReadAll(r.resp.Body)
}

func (r *ResponsePromise) setError(err error) {
	r.err = err
	r.ready <- struct{}{}
}

func (r *ResponsePromise) set(resp http.Response) {
	r.resp = resp
	r.ready <- struct{}{}
}

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

// FIXME: We should not have to define this here !
// TODO: Split transport.ResponseWriter vs protocol.ResponseWriter
func (w *HTTPResponseWriter) WriteError(err error) error {
	return nil
}

func (w *HTTPResponseWriter) Header() header.Header {
	return map[string][]string{} // FIXME: w.resp.Header()
}

type HTTPTransport struct {
	mux  *http.ServeMux
	reqs chan Data
}

// TODO: Pipelining !?
func NewHTTPTransport() *HTTPTransport {
	return &HTTPTransport{
		mux: http.NewServeMux(),
		// TODO: Buffered or unbuffered !?
		reqs: make(chan Data),
	}
}

func (t *HTTPTransport) Listen(address string) {
	t.mux.HandleFunc("/rpc/", t.handler)
	listner, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Listening failed: %s", err)
	}
	log.Printf("Listening on %s\n", listner.Addr())
	http.Serve(listner, t.mux)
}

func (t *HTTPTransport) handler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	// TODO: Handle errors.
	resp := HTTPResponseWriter{
		make(chan struct{}),
		w,
	}
	t.reqs <- Data{
		Body: body,
		Resp: &resp,
	}
	// Wait here else http.ResponseWriter became invalid.
	// TODO: Timeout !?
	<-resp.sent
}

func (t *HTTPTransport) Send(endpoint string, b []byte) (ResponseReader, error) {
	resp := ResponsePromise{ready: make(chan struct{})}
	go t.sendHTTP(endpoint, bytes.NewReader(b), &resp)
	// TODO: How about errors !?
	return &resp, nil
}

func (t *HTTPTransport) sendHTTP(endpoint string, body io.Reader, resp *ResponsePromise) {
	endpoint = fmt.Sprintf("%s/rpc/", endpoint)
	// TODO: content-type doesn't belong here.
	r, err := http.Post(endpoint, "application/json-rpc", body)
	if err != nil {
		resp.setError(err)
		return
	}
	// TODO: What should we when resp is not 200 !?
	resp.set(*r)
}

func (t *HTTPTransport) Receive() <-chan Data {
	return t.reqs
}
