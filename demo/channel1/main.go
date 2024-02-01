package main

import (
	"fmt"
	"time"
)

// ping-pong 模式
func main() {
	var ball int

	table := make(chan int)

	go player(table)
	go player(table)

	table <- ball
	time.Sleep(1 * time.Second)
	fmt.Println(<-table)
}

func player(table chan int) {
	for {
		ball := <-table
		ball++
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}
