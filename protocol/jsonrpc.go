package protocol

import (
	"fmt"

	"github.com/mouadino/go-nano/interfaces"
	"github.com/mouadino/go-nano/serializer"
)

type RequestBody struct {
	Version string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Id      string                 `json:"id"`
}

func (b *RequestBody) ToRequest() *interfaces.Request {
	return &interfaces.Request{
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
	interfaces.ResponseWriter
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
	fmt.Printf("Write %s\n", b)
	err = w.ResponseWriter.Write(b)
	if err != nil {
		return err
	}
	return nil
}

type JSONRPCProtocol struct {
	transport  interfaces.Transport
	serializer interfaces.Serializer
}

// TODO: DI for serializer !
func NewJSONRPCProtocol(t interfaces.Transport) *JSONRPCProtocol {
	return &JSONRPCProtocol{
		transport:  t,
		serializer: serializer.JSONSerializer{},
	}
}

// TODO: should we return Response !?
func (p *JSONRPCProtocol) SendRequest(endpoint string, r *interfaces.Request) (interface{}, error) {
	body := RequestBody{
		Version: "2.0",
		Method:  r.Method,
		Params:  r.Params,
		Id:      "0", // TODO: gouuid !?
	}
	b, err := p.serializer.Encode(body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Sending %v -> %s\n", b, endpoint)
	resp, err := p.transport.Send(endpoint, b)
	if err != nil {
		return nil, err
	}
	data, err := resp.Read()
	fmt.Printf("Received %s\n", data)
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

func (p *JSONRPCProtocol) ReceiveRequest() (interfaces.ResponseWriter, *interfaces.Request) {
	b := <-p.transport.Receive()
	body := RequestBody{}
	err := p.serializer.Decode(b.Body, &body)
	fmt.Printf("body %s %s\n", b.Body, body)
	if err != nil {
		return nil, nil
	}
	return &JSONRPCResponseWriter{
		b.Resp,
		p,
	}, body.ToRequest()
}
