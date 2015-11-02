package main

import (
	"log"
	"strings"
	"time"

	nano "github.com/mouadino/go-nano"
)

type echoService struct{}

func (echoService) NanoStart() error {
	log.Println("Starting ...")
	return nil
}

func (echoService) NanoStop() {
	log.Println("Stopping ...")
}

func (echoService) Upper(s string) string {
	time.Sleep(5 * time.Second)
	return strings.ToUpper(s)
}

func main() {
	nano.Default(echoService{}).ListenAndServe()
}
