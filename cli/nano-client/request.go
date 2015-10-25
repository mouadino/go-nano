package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mouadino/go-nano"
)

func SendRequest(c *cli.Context) {
	if c.Generic("service") == nil || c.Generic("params") == nil || c.Generic("method") == nil {
		fmt.Println("error: Missing argument, check help")
		os.Exit(1)
	}

	client := nano.Client(c.String("service"))

	var params map[string]interface{}
	err := json.Unmarshal([]byte(c.String("params")), &params)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	ret, err := client.Call(c.String("method"), params)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", ret)
}
