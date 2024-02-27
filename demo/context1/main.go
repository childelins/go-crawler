package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	before := time.Now()
	preCtx, _ := context.WithTimeout(ctx, 100*time.Millisecond)
	go func() {
		childCtx, _ := context.WithTimeout(preCtx, 300*time.Millisecond)
		select {
		case <-childCtx.Done():
			after := time.Now()
			fmt.Println("child during:", after.Sub(before).Milliseconds())
		}
	}()

	select {
	case <-preCtx.Done():
		after := time.Now()
		fmt.Println("pre during:", after.Sub(before).Milliseconds())
	}

	time.Sleep(1 * time.Second)

	// output:
	// child during: 100
	// pre during: 100

	// output (preCtx timeout = 500):
	// child during: 300
	// pre during: 500

	// 结论：父 Context 的退出会导致所有子 Context 的退出，而子 Context 的退出并不会影响父 Context。
}
