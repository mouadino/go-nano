package http

import "net/http"

// HTTPResponseWriter wraps around http ResponseWriter to implement
// go-nano transport.ResponseWriter.
type HTTPResponseWriter struct {
	sent chan struct{}
	resp http.ResponseWriter
}

// Write given data to underlying http ResponseWriter, data must
// be []byte.
func (w *HTTPResponseWriter) Write(data interface{}) error {
	_, err := w.resp.Write(data.([]byte))
	if err != nil {
		return err
	}
	w.sent <- struct{}{}
	return nil
}
