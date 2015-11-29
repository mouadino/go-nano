package jsonrpc

import (
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/transport"
)

// TODO: Rename message.
type ResponseBody struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *ErrorBody  `json:"error"`
	Id      string      `json:"id"`
}

func (r *ResponseBody) ToResponse() *protocol.Response {
	var err error
	if r.Error != nil {
		err = r.Error.Error()
	}
	return &protocol.Response{
		Body:  r.Result,
		Error: err,
		// TODO: Header
	}
}

func NewResponseBody(res interface{}, err error) *ResponseBody {
	return &ResponseBody{
		Version: "2.0",
		Result:  res,
		Error:   FromNanoError(err),
		Id:      "0", // TODO: Take same request id.
	}
}

type JSONRPCResponseWriter struct {
	transRW transport.ResponseWriter
	proto   *JSONRPCProtocol
	header  header.Header
}

func (rw *JSONRPCResponseWriter) Set(data interface{}) error {
	body := NewResponseBody(data, nil)
	return rw.writeToTransport(body)
}

func (rw *JSONRPCResponseWriter) SetError(err error) error {
	body := NewResponseBody(nil, err)
	return rw.writeToTransport(body)
}

func (rw *JSONRPCResponseWriter) writeToTransport(body *ResponseBody) error {
	b, err := rw.proto.serial.Encode(body)
	if err != nil {
		return err
	}
	// TODO: Write headers too !
	return rw.transRW.Write(b)
}

func (rw *JSONRPCResponseWriter) Header() header.Header {
	return rw.header
}
