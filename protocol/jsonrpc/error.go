package jsonrpc

import (
	"fmt"

	"github.com/mouadino/go-nano/protocol"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (b *ErrorBody) Error() error {
	return fmt.Errorf("jsonrpc error <%s> %s", b.Code, b.Message)
}

func FromNanoError(err error) *ErrorBody {
	// TODO: Set http status,
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
