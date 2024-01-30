package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "https://www.thepaper.cn/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("fetch url error:%v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%v", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read content failed:%v", err)
		return
	}

	// fmt.Println("body:", string(body))

	// TODO 功能是等价的，但是benchmark给出的结果是bytes库相较strings库性能会高很多

	// numLinks := strings.Count(string(body), "<a")
	numLinks := bytes.Count(body, []byte("<a"))
	fmt.Printf("homepage has %d links!\n", numLinks)

	// exist := strings.Contains(string(body), "疫情")
	exist := bytes.Contains(body, []byte("疫情"))
	fmt.Printf("是否存在疫情:%v\n", exist)
}
