package jsonrpc

import (
	"io"

	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/serializer"
)

type RequestBody struct {
	Version string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Id      string                 `json:"id"`
}

type ResponseBody struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *ErrorBody  `json:"error"`
	Id      string      `json:"id"`
}

func (b *RequestBody) ToRequest() *protocol.Request {
	return &protocol.Request{
		Method: b.Method,
		Params: b.Params,
		// TODO:
		Header: header.Header{},
	}
}

type jsonRPCProtocol struct {
	serial serializer.Serializer
}

func Serializer(serial serializer.Serializer) func(*jsonRPCProtocol) {
	return func(p *jsonRPCProtocol) {
		p.serial = serial
	}
}

func New(options ...func(*jsonRPCProtocol)) *jsonRPCProtocol {
	proto := &jsonRPCProtocol{
		serial: serializer.JSONSerializer{},
	}

	for _, opt := range options {
		opt(proto)
	}
	return proto
}

func (proto *jsonRPCProtocol) String() string {
	return "jsonrpc"
}

func (proto *jsonRPCProtocol) DecodeRequest(r io.Reader) (*protocol.Request, error) {
	req := RequestBody{}
	err := proto.serial.Decode(r, &req)
	if err != nil {
		return nil, err
	}
	// TODO: Headers ?
	return &protocol.Request{
		Method: req.Method,
		Params: req.Params,
	}, nil
}

func (proto *jsonRPCProtocol) EncodeRequest(r *protocol.Request) ([]byte, error) {
	body := RequestBody{
		Version: "2.0",
		Method:  r.Method,
		Params:  r.Params,
		Id:      "0", // TODO: gouuid !?
	}
	return proto.serial.Encode(body)
}

func (proto *jsonRPCProtocol) DecodeResponse(r io.Reader) (*protocol.Response, error) {
	resp := ResponseBody{}
	err := proto.serial.Decode(r, &resp)
	if err != nil {
		return nil, err
	}
	var Error error = nil
	if resp.Error != nil {
		Error = resp.Error.Error()
	}
	return &protocol.Response{
		Body:   resp.Result,
		Error:  Error,
		Header: header.Header{}, // TODO:
	}, nil
}

func (proto *jsonRPCProtocol) EncodeResponse(resp *protocol.Response) ([]byte, error) {
	body := ResponseBody{
		Version: "2.0",
		Result:  resp.Body,
		Error:   FromNanoError(resp.Error),
		Id:      "0", // TODO: Should match request id.
	}
	return proto.serial.Encode(body)
}
