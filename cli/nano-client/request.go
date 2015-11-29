package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mouadino/go-nano/client"
)

func SendRequest(c *cli.Context) {
	if c.Generic("service") == nil || c.Generic("params") == nil || c.Generic("method") == nil {
		fmt.Println("error: Missing argument, check help")
		os.Exit(1)
	}

	client := client.DefaultClient(c.String("service"))

	var params []interface{}
	err := json.Unmarshal([]byte(c.String("params")), &params)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	method := c.String("method")
	ret, err := client.Call(method, params...)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", ret)
}
