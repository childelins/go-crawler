package main

import (
	"fmt"
	"sync"
	"time"
)

type Stat struct {
	counters map[string]int64
	// 读写锁
	mutex sync.RWMutex
}

func (s *Stat) getCounter(name string) int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.counters[name]
}

func (s *Stat) SetCounter(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.counters[name]++
}

func main() {
	s := Stat{counters: map[string]int64{}}
	go s.SetCounter("hi")
	go s.SetCounter("hi")

	time.Sleep(1 * time.Second)
	fmt.Println(s.getCounter("hi"))
}
