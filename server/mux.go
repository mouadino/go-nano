package server

import (
	"fmt"
	"strings"

	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
)

type handlersMux struct {
	hdlrs map[string]handler.Handler
}

func (m *handlersMux) Names() []string {
	keys := make([]string, 0, len(m.hdlrs))
	for k := range m.hdlrs {
		keys = append(keys, k)
	}
	return keys
}

func (m *handlersMux) Add(name string, hdlr handler.Handler) error {
	_, exists := m.hdlrs[name]
	if exists {
		return fmt.Errorf("name %q already exist", name)
	}

	m.hdlrs[name] = hdlr
	return nil
}

func (m *handlersMux) Handle(resp *protocol.Response, req *protocol.Request) {
	parsedMethod := strings.SplitN(req.Method, ".", 2)
	hdlrName := parsedMethod[0]

	hdlr, ok := m.hdlrs[hdlrName]
	if !ok {
		// TODO: UnknownHandler error.
		resp.Error = fmt.Errorf("Unknown handler %q", hdlrName)
		return
	}
	if len(parsedMethod) > 1 {
		req.Method = parsedMethod[1]
	}
	hdlr.Handle(resp, req)
}
