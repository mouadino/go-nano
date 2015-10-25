package main

import (
	"fmt"
	"time"

	nano "github.com/mouadino/go-nano"
)

// FIXME: Using dynamic port now.
var echo = nano.Client("http://127.0.0.1:8080")

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
