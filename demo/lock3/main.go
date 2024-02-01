package main

import (
	"fmt"
	"sync"
	"time"
)

var count int64
var m sync.Mutex

// 互斥锁
func add() {
	m.Lock()
	count++
	m.Unlock()
}

func main() {
	go add()
	go add()

	time.Sleep(1 * time.Second)
	fmt.Println(count)
}
