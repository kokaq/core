package main

import (
	"fmt"

	"github.com/kokaq/core/queue"
)

func main() {
	queue, err := queue.NewDefaultKokaq(1, 1)
	if err != nil {
		panic(err)
	}
	if queue.IsEmpty() {
		fmt.Println("Is Empty")
	}
}
