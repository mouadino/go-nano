package middleware

import (
	"bytes"
	"regexp"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mouadino/go-nano/header"
	"github.com/mouadino/go-nano/protocol"
)

func TestLoggerMiddleware(t *testing.T) {
	req := &protocol.Request{
		Method: "foobar",
		Params: protocol.Params{},
		Header: header.Header{},
	}
	resp := &protocol.Response{}

	buff := bytes.NewBufferString("")
	logger := log.New()
	logger.Out = buff

	handler := Chain(&dummyHandler{}, NewLogger(logger))

	handler.Handle(resp, req)

	loglines := buff.String()

	expectedLog := regexp.MustCompile("method=foobar duration=[0-9]+")
	if expectedLog.FindString(loglines) != "" {
		t.Errorf("didn't find expected log line in %s", loglines)
	}
}
