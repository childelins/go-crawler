package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// 反向代理
func main() {
	// 初始化反向代理服务
	proxy, err := NewProxy()
	if err != nil {
		panic(err)
	}

	// 所有请求都由ProxyRequestHandler函数进行处理
	http.HandleFunc("/", ProxyRequestHandler(proxy))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ProxyRequestHandler 使用代理处理HTTP请求
func ProxyRequestHandler(proxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func NewProxy() (*httputil.ReverseProxy, error) {
	// 实际的后端服务器地址
	targetHost := "http://my-api-server.com"
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	return proxy, nil
}
