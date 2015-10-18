package transport

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mouadino/go-nano/interfaces"
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

type ResponseWriter struct {
	sent chan struct{}
	resp http.ResponseWriter
}

func (w *ResponseWriter) Write(data interface{}) error {
	_, err := w.resp.Write(data.([]byte))
	if err != nil {
		fmt.Printf("Error %s\n", err)
		return err
	}
	w.sent <- struct{}{}
	return nil
}

func (w *ResponseWriter) Header() interfaces.Header {
	return map[string][]string{} // FIXME: w.resp.Header()
}

type HTTPTransport struct {
	reqs chan interfaces.Data
}

// TODO: Pipelining !?
func NewHTTPTransport(address string) *HTTPTransport {
	t := &HTTPTransport{
		// TODO: Buffered or unbuffered !?
		reqs: make(chan interfaces.Data),
	}
	http.HandleFunc("/rpc/", t.handler)
	// TODO: Start !?
	go http.ListenAndServe(address, nil)
	return t
}

func (t *HTTPTransport) handler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	// TODO: Handle errors.
	resp := ResponseWriter{
		make(chan struct{}),
		w,
	}
	t.reqs <- interfaces.Data{
		Body: body,
		Resp: &resp,
	}
	// Wait here else ResponseWriter became invalid.
	// TODO: Timeout !?
	<-resp.sent
}

func (t *HTTPTransport) Send(endpoint string, b []byte) (interfaces.ResponseReader, error) {
	resp := ResponsePromise{}
	go t.sendHTTP(endpoint, bytes.NewReader(b), &resp)
	// TODO: How about errors !?
	return &resp, nil
}

func (t *HTTPTransport) sendHTTP(endpoint string, body io.Reader, resp *ResponsePromise) {
	r, err := http.Post(endpoint, "application/json-rpc", body)
	if err != nil {
		resp.setError(err)
		return
	}
	resp.set(*r)
}

func (t *HTTPTransport) Receive() <-chan interfaces.Data {
	return t.reqs
}
