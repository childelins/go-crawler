package main

import (
	"fmt"
	"sync"
	"time"
)

var once sync.Once

func DoOnce() {
	once.Do(func() {
		fmt.Println("Am called")
	})
}

func main() {
	go DoOnce()
	go DoOnce()

	time.Sleep(1 * time.Second)
}
