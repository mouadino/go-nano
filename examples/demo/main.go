package main

import (
	"fmt"
	"time"

	"github.com/mouadino/go-nano/client"
)

var echo = client.DefaultClient("jsonrpc+amqp://127.0.0.1:5672/upper")

/*var amqpTrans = transport.NewAMQPTransport("amqp://127.0.0.1:5672")
var echo = client.CustomClient(
	"upper",
	jsonrpc.NewJSONRPCProtocol(amqpTrans, serializer.JSONSerializer{}),
	extension.NewTimeoutExt(3*time.Second),
)*/

func main() {
	c := time.Tick(1 * time.Second)
	i := 0
	for _ = range c {
		text := fmt.Sprintf("foo_%d", i)
		result, err := echo.Call("upper", text)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("%s\n", result.(string))
		}
		i++
	}
}
