package main

import (
	"fmt"
	"go-crawler/collect"
	"go-crawler/proxy"
	"time"
)

func main() {
	url := "https://book.douban.com/subject/1007305/"

	// error status code:418
	// f := collect.BaseFetch{}

	proxyUrls := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	p, err := proxy.RoundRobinProxySwitcher(proxyUrls...)
	if err != nil {
		panic(err)
	}

	f := collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}
	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("read content failed: %v", err)
		return
	}

	fmt.Println(string(body))
}
