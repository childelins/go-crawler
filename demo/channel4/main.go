package main

import (
	"fmt"
	"sync"
	"time"
)

// fan-out 模式
func main() {
	var wg sync.WaitGroup
	wg.Add(36)
	go pool(&wg, 36, 50)
	wg.Wait()
}

func pool(wg *sync.WaitGroup, workers, tasks int) {
	tashCh := make(chan int)

	for i := 0; i < workers; i++ {
		go worker(tashCh, wg)
	}

	for i := 0; i < tasks; i++ {
		tashCh <- i
	}

	close(tashCh)
}

func worker(tashCh <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, ok := <-tashCh
		if !ok {
			return
		}

		d := time.Duration(task) * time.Millisecond
		time.Sleep(d)
		fmt.Println("processing task", task)
	}
}
