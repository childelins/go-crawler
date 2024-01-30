package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func main() {
	// 访问路由到hello函数
	http.HandleFunc("/", hello)
	// 监听本地8080端口
	http.ListenAndServe(":8080", nil)
}
