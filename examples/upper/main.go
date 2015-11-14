package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"

	nano "github.com/mouadino/go-nano"
	"github.com/mouadino/go-nano/discovery"
)

type echoService struct{}

func (echoService) NanoStart() error {
	log.Debug("Starting ...")
	return nil
}

func (echoService) NanoStop() {
	log.Debug("Stopping ...")
}

func (echoService) Upper(s string) string {
	return strings.ToUpper(s)
}

func main() {
	zkAnnouncer := discovery.DefaultZooKeeperAnnounceResolver(
		[]string{"127.0.0.1:2181"},
	)
	server := nano.DefaultServer(echoService{})
	server.Announce("upper", discovery.ServiceMetadata{}, zkAnnouncer)
	server.ListenAndServe()
}
