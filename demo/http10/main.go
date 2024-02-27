package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

// HTTP 隧道代理
func main() {
	server := &http.Server{
		Addr: ":9981",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
	}

	log.Fatal(server.ListenAndServe())
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	// 与目标服务器建立TCP连接
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	// 通过 hijacker.Hijack() 拿到了客户端与代理服务器之间的底层 TCP 连接
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()

	// 通过 io.Copy 就简单地串联起了一个管道，实现了数据包在服务器与客户端之间的相互转发
	// 当然在工业级的代码中，我们不会用这么粗暴的方式实现这一功能，因为传输的数据量可能很大
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
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
