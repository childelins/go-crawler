package main

import (
	"fmt"
	"go-crawler/collect"
	"go-crawler/engine"
	"go-crawler/log"
	"go-crawler/parse/doubangroup"
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
	// =============== test =====================
	// url := "https://book.douban.com/subject/1007305/"
	// body, err := f.Get(url)
	// if err != nil {
	// 	fmt.Printf("read content failed: %v", err)
	// 	return
	// }

	// logger.Info("get content", zap.Int("len", len(body)))

	// // proxy
	// proxyUrls := []string{"http://127.0.0.1:8888", "http://127.0.0.1:8889"}
	// p, err := proxy.RoundRobinProxySwitcher(proxyUrls...)
	// if err != nil {
	// 	logger.Error("RoundRobinProxySwitcher failed")
	// 	return
	// }

	// url := "<https://google.com>"
	// f := collect.BrowserFetch{
	// 	Timeout: 3000 * time.Millisecond,
	// 	Proxy:   p,
	// }

	// body, err := f.Get(url)
	// if err != nil {
	// 	fmt.Printf("read content failed:%v\\n", err)
	// 	return
	// }

	// fmt.Println(string(body))

	// douban cookie
	// cookie := `bid=qk-KbS-ffCg; douban-fav-remind=1; Hm_lvt_6d4a8cfea88fa457c3127e14fb5fabc2=1700728840; _ga=GA1.2.478668193.1700728840; ll="118201"; _ga_Y4GN1R87RG=GS1.1.1700728839.1.1.1700728894.0.0.0; viewed="1007305_35196328_35474931_35219951_36449803_36368057_36424128"; ap_v=0,6.0; __utmc=30149280; __utmz=30149280.1706838836.7.6.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/column/article/612328; __utma=30149280.125988448.1699844939.1706838836.1706844162.8; __utmt=1; dbcl2="126963156:6rHEcFsnXwM"; ck=Focf; push_noty_num=0; push_doumail_num=0; ct=y;`

	// v1
	// var worklist []*collect.Request
	// for i := 0; i <= 100; i += 25 {
	// 	url := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
	// 	worklist = append(worklist, &collect.Request{
	// 		Url:       url,
	// 		Cookie:    cookie,
	// 		ParseFunc: doubangroup.ParseURL,
	// 	})
	// }

	// f := collect.BrowserFetch{
	// 	Timeout: 3000 * time.Millisecond,
	// 	Proxy:   nil,
	// }

	// for len(worklist) > 0 {
	// 	items := worklist
	// 	worklist = nil
	// 	for _, item := range items {
	// 		body, err := f.Get(item)
	// 		// 休眠 1 秒钟尽量减缓服务器的压力
	// 		time.Sleep(1 * time.Second)
	// 		if err != nil {
	// 			logger.Error("read content failed", zap.Error(err))
	// 			continue
	// 		}

	// 		res := item.ParseFunc(body, item)
	// 		for _, item := range res.Items {
	// 			logger.Info("result", zap.String("get url:", item.(string)))
	// 		}

	// 		// 广度优先搜索
	// 		worklist = append(worklist, res.Requests...)
	// 	}
	// }

	// v2
	seeds := make([]*collect.Task, 0, 1000)
	cookie := `bid=qk-KbS-ffCg; douban-fav-remind=1; Hm_lvt_6d4a8cfea88fa457c3127e14fb5fabc2=1700728840; _ga=GA1.2.478668193.1700728840; ll="118201"; _ga_Y4GN1R87RG=GS1.1.1700728839.1.1.1700728894.0.0.0; viewed="1007305_35196328_35474931_35219951_36449803_36368057_36424128"; ap_v=0,6.0; __utmc=30149280; __utmz=30149280.1706838836.7.6.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/column/article/612328; __utma=30149280.125988448.1699844939.1706838836.1706844162.8; __utmt=1; dbcl2="126963156:6rHEcFsnXwM"; ck=Focf; push_noty_num=0; push_doumail_num=0; ct=y;`

	for i := 0; i <= 100; i += 25 {
		url := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		seeds = append(seeds, &collect.Task{
			Url:      url,
			Cookie:   cookie,
			WaitTime: 1 * time.Second,
			MaxDepth: 5,
			RootReq: &collect.Request{
				ParseFunc: doubangroup.ParseURL,
			},
		})
	}

	f := collect.BrowserFetch{
		Timeout: 5000 * time.Millisecond,
		Proxy:   nil,
		Logger:  logger,
	}

	s := engine.NewSchedule(
		engine.WithFetcher(f),
		engine.WithLogger(logger),
		engine.WithWorkCount(5),
		engine.WithSeeds(seeds),
	)

	s.Run()
}
