package main

import (
	"fmt"
	"go-crawler/collect"
	"time"
)

func main() {
	url := "https://book.douban.com/subject/1007305/"

	// error status code:418
	// f := collect.BaseFetch{}
	f := collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
	}
	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("read content failed: %v", err)
		return
	}

	fmt.Println(string(body))
}
