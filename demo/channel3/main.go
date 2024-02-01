package main

import (
	"fmt"
	"time"
)

// fan-in 模式 -> 多路复用
func search(msg string) chan string {
	ch := make(chan string)

	go func() {
		var i int
		for {
			// 模拟找到了关键字
			ch <- fmt.Sprintf("get %s %d", msg, i)
			i++
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	return ch
}

func main() {
	ch1 := search("jonson")
	ch2 := search("olaya")

	for {
		select {
		case msg := <-ch1:
			fmt.Println(msg)
		case msg := <-ch2:
			fmt.Println(msg)
		}
	}
}
