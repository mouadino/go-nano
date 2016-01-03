package protocol

import "errors"

var (
	UnknownMethod = errors.New("unknown method")
	ParamsError   = errors.New("unknown parameters")
	InternalError = errors.New("Internal error")
	ServerError   = errors.New("Server error")
)
