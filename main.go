package main

import (
	"fmt"
	"go-crawler/collect"
	"go-crawler/log"
	"go-crawler/proxy"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	// plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	// defer c.Close()

	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger = log.NewLogger(plugin)
	defer logger.Sync()
}

func main() {
	logger.Debug("log init end")

	url := "https://book.douban.com/subject/1007305/"

	// error status code:418
	// f := collect.BaseFetch{}

	proxyUrls := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	p, err := proxy.RoundRobinProxySwitcher(proxyUrls...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
		return
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

	logger.Info("get content", zap.Int("len", len(body)))
}
