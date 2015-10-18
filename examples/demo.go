// FIXME: No Implementation yet!
package main

import (
	"fmt"
	"time"

	nano "github.com/mouadino/go-nano"
)

var echo = nano.Client("echo")

type demoService struct {
	Delay int
}

// TODO: Do I need to make it a service ?
// TODO: How to start Main ?
func (svc demoService) Main(s string) string {
	c := time.Tick(svc.Delay * time.Second)
	for i, _ := range c {
		text := fmt.Sprintf("foo_%s", i)
		result := echo.Echo(text)
		fmt.Println("%v", result)
	}
}

func main() {
	// TODO: Delay as config.
	nano.Main(demoService{
		Delay: 1,
	})
}
