package main

import (
	"fmt"
	"time"

	nano "github.com/mouadino/go-nano"
)

var echo = nano.Client("http://127.0.0.1:8080")

func main() {
	c := time.Tick(1 * time.Second)
	i := 0
	for _ = range c {
		fmt.Println("Calling ...")
		text := fmt.Sprintf("foo_%d", i)
		result, err := echo.Call("upper", map[string]interface{}{"text": text})
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Printf("%s\n", result.(string))
		}
		i++
	}
}
