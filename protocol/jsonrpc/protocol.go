package jsonrpc

import (
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/serializer"
	"github.com/mouadino/go-nano/transport"
)

type RequestBody struct {
	Version string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Id      string                 `json:"id"`
}

func (b *RequestBody) ToRequest() *protocol.Request {
	return &protocol.Request{
		Method: b.Method,
		Params: b.Params,
	}
}

type JSONRPCProtocol struct {
	trans  transport.Transport
	serial serializer.Serializer
}

func NewJSONRPCProtocol(trans transport.Transport, serial serializer.Serializer) *JSONRPCProtocol {
	return &JSONRPCProtocol{
		trans:  trans,
		serial: serial,
	}
}

func (proto *JSONRPCProtocol) SendRequest(endpoint string, r *protocol.Request) (interface{}, error) {
	reqBody, err := proto.getBody(r)
	if err != nil {
		return nil, err
	}
	resp, err := proto.trans.Send(endpoint, reqBody)
	respBody := ResponseBody{}
	err = proto.serial.Decode(resp, &respBody)
	if err != nil {
		return nil, err
	}
	if respBody.Error != nil {
		return nil, respBody.Error.Error()
	}
	return respBody.Result, nil
}

func (proto *JSONRPCProtocol) getBody(r *protocol.Request) ([]byte, error) {
	body := RequestBody{
		Version: "2.0",
		Method:  r.Method,
		Params:  r.Params,
		Id:      "0", // TODO: gouuid !?
	}
	return proto.serial.Encode(body)
}

func (proto *JSONRPCProtocol) ReceiveRequest() (protocol.ResponseWriter, *protocol.Request) {
	b := <-proto.trans.Receive()
	body := RequestBody{}
	err := proto.serial.Decode(b.Body, &body)
	if err != nil {
		return nil, nil
	}
	rw := &JSONRPCResponseWriter{
		b.Resp,
		proto,
		header.Header{},
	}
	return rw, body.ToRequest()
}
