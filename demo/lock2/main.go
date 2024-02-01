package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var flag int64
var count int64

// 自旋锁
func add() {
	for {
		if atomic.CompareAndSwapInt64(&flag, 0, 1) {
			count++
			atomic.StoreInt64(&flag, 0)
			return
		}
	}
}

func main() {
	go add()
	go add()

	time.Sleep(1 * time.Second)
	fmt.Println(atomic.LoadInt64(&count))
}
