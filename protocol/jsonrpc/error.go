package jsonrpc

import "github.com/mouadino/go-nano/protocol"

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (b *ErrorBody) Error() error {
	switch b.Code {
	case "-32601":
		return protocol.UnknownMethod
	case "-32602":
		return protocol.ParamsError
	case "-32603":
		return protocol.InternalError
	default:
		return protocol.ServerError
	}
}

func FromNanoError(err error) *ErrorBody {
	if err == nil {
		return nil
	}
	switch {
	case err == protocol.UnknownMethod:
		return &ErrorBody{
			Code:    "-32601",
			Message: err.Error(),
			Data:    "",
		}
	case err == protocol.ParamsError:
		return &ErrorBody{
			Code:    "-32602",
			Message: err.Error(),
			Data:    "",
		}
	case err == protocol.InternalError:
		return &ErrorBody{
			Code:    "-32603",
			Message: err.Error(),
			Data:    "",
		}
	default:
		return &ErrorBody{
			Code:    "-32000",
			Message: "Server error",
			Data:    err.Error(),
		}
	}
}
