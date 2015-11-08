package middleware

import (
	"bytes"
	"regexp"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
)

var panicHandler = handler.HandlerFunc(func(rw protocol.ResponseWriter, req *protocol.Request) {
	panic("Booom")
})

func TestRecoverMiddleware(t *testing.T) {
	req := &protocol.Request{
		Method: "foobar",
		Params: protocol.Params{},
		Header: header.Header{},
	}
	rw := &protocol.DumpResponseWriter{
		HeaderValues: header.Header{},
	}

	buff := bytes.NewBufferString("")
	logger := log.New()
	logger.Out = buff

	handler := Chain(panicHandler, NewRecoverMiddleware(logger, true, 8*1024))

	handler.Handle(rw, req)

	loglines := buff.String()

	expectedLog := regexp.MustCompile(`error="Booom"`)
	if expectedLog.FindString(loglines) != "" {
		t.Errorf("didn't find expected log line in %s", loglines)
	}
}
