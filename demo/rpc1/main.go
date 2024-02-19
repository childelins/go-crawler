package main

import (
	"go-crawler/demo/rpc1/server"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	// 注册rpc服务
	rpc.Register(new(server.Arith))
	// 采用http协议作为rpc载体
	rpc.HandleHTTP()

	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	// 常规启动http服务
	go http.Serve(ln, nil)

	done := make(chan struct{})
	<-done
}
