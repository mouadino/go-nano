package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/transport"
	"github.com/mouadino/go-nano/utils"
)

const RPCPath = "/rpc/"

type HTTPTransport struct {
	mux       *http.ServeMux
	reqs      chan transport.Request
	addr      string
	listening bool
}

// New creates a new HTTP transport.
func New() *HTTPTransport {
	return &HTTPTransport{
		mux:  http.NewServeMux(),
		reqs: make(chan transport.Request),
	}
}

// Listen instruct the http server to listen on random port and
// external ip of the node.
func (trans *HTTPTransport) Listen() error {
	trans.mux.HandleFunc(RPCPath, trans.handle)
	ip, err := utils.GetExternalIP()
	if err != nil {
		return err
	}
	listner, err := net.Listen("tcp", fmt.Sprintf("%s:0", ip))
	if err != nil {
		return err
	}
	trans.addr = fmt.Sprintf("http://%s", listner.Addr().String())
	log.Info("Listening on ", trans.addr)
	go http.Serve(listner, trans.mux)
	trans.listening = true
	return nil
}

// Addr returns listening address in the form: http://...:...
func (trans *HTTPTransport) Addr() string {
	return trans.addr
}

func (trans *HTTPTransport) handle(rw http.ResponseWriter, req *http.Request) {
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
	// TODO: Timeout !?, http.Hijacker ? Can we do better ?
	<-resp.sent
}

// Send a raw request to given endpoint.
func (trans *HTTPTransport) Send(endpoint string, body io.Reader) ([]byte, error) {
	endpoint += RPCPath
	// TODO: content-type doesn't belong here.
	// TODO: Send RPC headers as HTTP headers ?
	resp, err := http.Post(endpoint, "application/json-rpc", body)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// Receive returns received requests.
func (trans *HTTPTransport) Receive() <-chan transport.Request {
	return trans.reqs
}
