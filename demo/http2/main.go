package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// http请求
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	// 获取返回的数据
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 打印返回的数据
	fmt.Println(string(content))
}
