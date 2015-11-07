package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"

	nano "github.com/mouadino/go-nano"
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
	nano.Default(echoService{}).ListenAndServe()
}
