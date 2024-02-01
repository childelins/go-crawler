package main

import (
	"fmt"
	"time"
)

// fan-in 模式
func search(ch chan string, msg string) {
	var i int
	for {
		// 模拟找到了关键字
		ch <- fmt.Sprintf("get %s %d", msg, i)
		i++
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	ch := make(chan string)
	go search(ch, "jonson")
	go search(ch, "olaya")

	for i := range ch {
		fmt.Println(i)
	}
}
