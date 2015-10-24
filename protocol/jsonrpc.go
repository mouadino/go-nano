package protocol

import (
	"log"

	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
)

type RequestBody struct {
	Version string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Id      string                 `json:"id"`
}

func (b *RequestBody) ToRequest() *Request {
	return &Request{
		Method: b.Method,
		Params: b.Params,
	}
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type ResponseBody struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   ErrorBody   `json:"error"`
	Id      string      `json:"id"`
}

// TODO: Use me !
func (b *ResponseBody) ToRespone() {
}

type JSONRPCResponseWriter struct {
	transport.ResponseWriter
	p *JSONRPCProtocol
}

func (w *JSONRPCResponseWriter) Write(data interface{}) error {
	body := ResponseBody{
		Version: "2.0",
		Result:  data,
		Id:      "0",
	}
	b, err := w.p.serializer.Encode(body)
	if err != nil {
		return err
	}
	log.Printf("Write %s\n", b)
	err = w.ResponseWriter.Write(b)
	if err != nil {
		return err
	}
	return nil
}

type JSONRPCProtocol struct {
	transport  transport.Transport
	serializer serializer.Serializer
}

// TODO: DI for serializer !
func NewJSONRPCProtocol(t transport.Transport) *JSONRPCProtocol {
	return &JSONRPCProtocol{
		transport:  t,
		serializer: serializer.JSONSerializer{},
	}
}

// TODO: should we return Response !?
func (p *JSONRPCProtocol) SendRequest(endpoint string, r *Request) (interface{}, error) {
	b, err := p.getBody(r)
	if err != nil {
		return nil, err
	}
	log.Printf("Sending %v -> %s\n", b, endpoint)
	resp, err := p.transport.Send(endpoint, b)
	if err != nil {
		return nil, err
	}
	data, err := resp.Read()
	log.Printf("Received %s\n", data)
	if err != nil {
		return nil, err
	}
	respBody := ResponseBody{}
	err = p.serializer.Decode(data, &respBody)
	if err != nil {
		return nil, err
	}
	// TODO: If body.Error return error
	return respBody.Result, nil
}

func (p *JSONRPCProtocol) getBody(r *Request) ([]byte, error) {
	body := RequestBody{
		Version: "2.0",
		Method:  r.Method,
		Params:  r.Params,
		Id:      "0", // TODO: gouuid !?
	}
	return p.serializer.Encode(body)
}

func (p *JSONRPCProtocol) ReceiveRequest() (transport.ResponseWriter, *Request) {
	b := <-p.transport.Receive()
	body := RequestBody{}
	err := p.serializer.Decode(b.Body, &body)
	log.Printf("body %s %s\n", b.Body, body)
	if err != nil {
		return nil, nil
	}
	return &JSONRPCResponseWriter{
		b.Resp,
		p,
	}, body.ToRequest()
}
