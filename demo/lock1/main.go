package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var count int64

// 原子锁
func add() {
	atomic.AddInt64(&count, 1)
}

func main() {
	go add()
	go add()

	time.Sleep(1 * time.Second)
	fmt.Println(atomic.LoadInt64(&count))
}
