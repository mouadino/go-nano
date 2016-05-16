package http

import (
	"net/http"

	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
)

// HTTPResponseWriter wraps around http ResponseWriter to implement
// go-nano transport.ResponseWriter.
type HTTPResponseWriter struct {
	proto protocol.Protocol
	resp  http.ResponseWriter
}

// Write given data to underlying http ResponseWriter, data must
// be []byte.
func (w *HTTPResponseWriter) Set(data interface{}) error {
	// TODO: Encode using proto
	_, err := w.resp.Write(data.([]byte))
	return err
}

func (w *HTTPResponseWriter) SetError(err error) error {
	// TODO
	return nil
}

func (w *HTTPResponseWriter) Header(err error) header.Header {
	// TODO
	return header.Header{}
}
