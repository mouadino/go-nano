package http

import (
	"errors"
	"net/http"
	"sync"

	"github.com/mouadino/go-nano/protocol"
)

type protocolMap struct {
	mu     sync.RWMutex
	protos map[string]protocol.Protocol
}

func (pm *protocolMap) Get(req *http.Request) protocol.Protocol {
	contentType := req.Header.Get("Content-Type")
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	if p, ok := pm.protos[contentType]; ok {
		return p
	}
	return nil
}

func (pm *protocolMap) Add(proto protocol.Protocol) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	name := proto.String()
	if _, ok := pm.protos[name]; ok {
		return errors.New("already exists")
	}
	pm.protos[name] = proto
	return nil
}
