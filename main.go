package main

import (
	"fmt"
	"time"

	"github.com/deharahawa/mba-go/multithreading/challenge/src/api"
)

func main() {
	ch := make(chan string)

	go api.CallBrasilAPI("01153000", ch)
	go api.CallViaCEPAPI("01153000", ch)

	select {
	case response := <-ch:
		fmt.Println(response)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")
	}
}
