package main

import (
	"errors"
	"strings"

	"github.com/mouadino/go-nano"
	"github.com/mouadino/go-nano/discovery"
	"github.com/mouadino/go-nano/handler"
	"github.com/mouadino/go-nano/protocol"
)

func Upper(rw protocol.ResponseWriter, req *protocol.Request) {
	text, ok := req.Params["_0"]
	if !ok {
		rw.WriteError(errors.New("Not ok"))
		return
	}
	rw.Write(strings.ToUpper(text.(string)))
}

func main() {
	zkAnnouncer := discovery.DefaultZooKeeperAnnounceResolver(
		[]string{"127.0.0.1:2181"},
	)
	server := nano.DefaultServer(handler.HandlerFunc(Upper))
	server.Announce("upper", discovery.ServiceMetadata{}, zkAnnouncer)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
