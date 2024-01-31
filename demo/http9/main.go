package main

import (
	"io"
	"log"
	"net/http"
)

// 正向代理
func main() {
	server := &http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleHTTP(w, r)
		}),
	}

	log.Fatal(server.ListenAndServe())
}

// 当前代理服务器获取客户端的请求，并用自己的身份发送请求到服务器。
// 代理服务器获取到服务器的回复后，会再次利用 io.Copy 将回复发送回客户端
func handleHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("=====request comming =====")
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Println("proxy err:", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
