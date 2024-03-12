package main

import (
	"go-crawler/collect"
	"go-crawler/engine"
	"go-crawler/log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	// plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	// defer c.Close()

	// log
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger = log.NewLogger(plugin)
	defer logger.Sync()

	logger.Info("log init end")
}

func main() {
	// // proxy
	// proxyURLs := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8888"}
	// p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	// if err != nil {
	// 	logger.Error("RoundRobinProxySwitcher failed")
	// }

	f := collect.BrowserFetch{
		Timeout: 5000 * time.Millisecond,
		Proxy:   nil,
		Logger:  logger,
	}

	seeds := make([]*collect.Task, 0, 1000)
	seeds = append(seeds, &collect.Task{
		Property: collect.Property{
			// Name: "find_douban_sun_room",
			Name: "douban_book_list",
		},
	})

	e := engine.NewEngine(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
		engine.WithScheduler(engine.NewSchedule()),
	)

	e.Run()
}
