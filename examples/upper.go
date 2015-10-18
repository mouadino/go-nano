package main

import (
	"strings"

	nano "github.com/mouadino/go-nano"
)

/* TODO: Do I need this ?
type EchoService interface {
	Echo(string) string
	Upper(string) string
}*/

type echoService struct{}

func (echoService) Echo(s string) string {
	return s
}

func (echoService) Upper(s string) string {
	return strings.ToUpper(s)
}

func main() {
	nano.Main(echoService{})
}
