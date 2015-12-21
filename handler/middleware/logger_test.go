package middleware

import (
	"bytes"
	"regexp"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
	"github.com/mouadino/go-nano/protocol/dummy"
)

func TestLoggerMiddleware(t *testing.T) {
	req := &protocol.Request{
		Method: "foobar",
		Params: protocol.Params{},
		Header: header.Header{},
	}
	rw := &dummy.ResponseWriter{
		HeaderValues: header.Header{},
	}

	buff := bytes.NewBufferString("")
	logger := log.New()
	logger.Out = buff

	handler := Chain(&dummyHandler{}, NewLogger(logger))

	handler.Handle(rw, req)

	loglines := buff.String()

	expectedLog := regexp.MustCompile("method=foobar duration=[0-9]+")
	if expectedLog.FindString(loglines) != "" {
		t.Errorf("didn't find expected log line in %s", loglines)
	}
}
