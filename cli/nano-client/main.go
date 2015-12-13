package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "nano-client"
	app.Usage = "Send a request to service"
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "service, s",
			Usage: "Service endpoint to send request to (Required)",
		},
		cli.StringFlag{
			Name:  "method, m",
			Usage: "RPC method to call (Required)",
		},
		cli.StringFlag{
			Name:  "params, p",
			Usage: "Parameters as JSON (Required)",
		},
	}
	//app.Action = SendRequest

	app.Run(os.Args)
}
