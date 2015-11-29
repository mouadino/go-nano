package protocol

import "github.com/mouadino/go-nano/header"

type DumpResponseWriter struct {
	Data         interface{}
	Error        error
	HeaderValues header.Header
}

func (rw *DumpResponseWriter) Header() header.Header {
	return rw.HeaderValues
}

func (rw *DumpResponseWriter) Set(data interface{}) error {
	rw.Data = data
	return nil
}

func (rw *DumpResponseWriter) SetError(err error) error {
	rw.Error = err
	return nil
}
