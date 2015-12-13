package http

import "net/http"

type HTTPResponseWriter struct {
	sent chan struct{}
	resp http.ResponseWriter
}

func (w *HTTPResponseWriter) Write(data interface{}) error {
	_, err := w.resp.Write(data.([]byte))
	if err != nil {
		return err
	}
	w.sent <- struct{}{}
	return nil
}
