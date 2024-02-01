package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	users []string
	cond  *sync.Cond
)

func givePrizes() {
	cond.L.Lock()
	for len(users) < 5 {
		// 循环等待
		cond.Wait()
	}

	fmt.Println("start give prize to users: ", users[:3])
	cond.L.Unlock()
}

func newUser(name string) {
	cond.L.Lock()
	users = append(users, name)
	cond.L.Unlock()

	// 唤醒一个最先等待的协程
	cond.Signal()
}

func main() {
	cond = sync.NewCond(&sync.Mutex{})

	go givePrizes()

	newUser("Lisi")
	newUser("Zhangsan")
	newUser("Wangwu")
	newUser("Zhaoliu")
	newUser("ZhouJun")

	time.Sleep(1 * time.Second)
}
