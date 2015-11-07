package jsonrpc

import "github.com/mouadino/go-nano/transport"

type ResponseBody struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *ErrorBody  `json:"error"`
	Id      string      `json:"id"`
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
}

func (rw *JSONRPCResponseWriter) Write(data interface{}) error {
	body := NewResponseBody(data, nil)
	return rw.writeToTransport(body)
}

func (rw *JSONRPCResponseWriter) WriteError(err error) error {
	body := NewResponseBody(nil, err)
	return rw.writeToTransport(body)
}

func (rw *JSONRPCResponseWriter) writeToTransport(body *ResponseBody) error {
	b, err := rw.proto.serial.Encode(body)
	if err != nil {
		return err
	}
	return rw.transRW.Write(b)
}
