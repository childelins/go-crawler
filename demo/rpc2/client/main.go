package main

import (
	"fmt"
	"log"
	"net/rpc/jsonrpc"
)

func main() {
	// 只有这里不一样
	client, err := jsonrpc.Dial("tcp", ":6001")
	if err != nil {
		log.Fatalln("dialing error:", err)
	}

	var result string
	err = client.Call("App.Hi", "jsonrpc", &result)
	if err != nil {
		log.Fatalln("call error:", err)
	}

	fmt.Printf("result = %s\n", result)
}
