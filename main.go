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
	// cookie := `bid=qk-KbS-ffCg; douban-fav-remind=1; _pk_id.100001.8cb4=2b9cc6a45e496b68.1699844938.; ll="118201"; viewed="1007305_35196328_35474931_35219951_36449803_36368057_36424128"; __utmz=30149280.1706838836.7.6.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/column/article/612328; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12696; _ga=GA1.2.478668193.1700728840; _ga_Y4GN1R87RG=GS1.1.1706844569.2.1.1706844957.0.0.0; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1709532615%2C%22https%3A%2F%2Ftime.geekbang.org%2Fcolumn%2Farticle%2F612328%22%5D; _pk_ses.100001.8cb4=1; ap_v=0,6.0; __utma=30149280.125988448.1699844939.1709286691.1709532615.11; __utmc=30149280; __utmt=1; dbcl2="126963156:PCmGfH4AjBA"; ck=lbgb; __utmb=30149280.29.5.1709534687162`

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
		Name: "find_douban_sun_room",
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
