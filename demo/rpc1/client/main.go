package main

import (
	"fmt"
	"go-crawler/demo/rpc1/server"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := server.Args{A: 7, B: 8}
	var reply int
	// 同步调用
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}

	fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

	quotient := new(server.Quotient)
	// 异步调用
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	<-divCall.Done

	fmt.Printf("Arith: %d/%d=%d\n", args.A, args.B, quotient)
}
