package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type App struct{}

func (s *App) Hi(name string, r *string) error {
	*r = fmt.Sprintf("Hello, %s!", name)
	return nil
}

func main() {
	rpc := rpc.NewServer()
	app := new(App)
	rpc.Register(app)

	ln, err := net.Listen("tcp", ":6001")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		// 使用 jsonrpc 作为解码器
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
