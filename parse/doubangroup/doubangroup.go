package doubangroup

import (
	"fmt"
	"go-crawler/collect"
	"regexp"
	"time"
)

const cookie = `bid=qk-KbS-ffCg; douban-fav-remind=1; _pk_id.100001.8cb4=2b9cc6a45e496b68.1699844938.; ll="118201"; viewed="1007305_35196328_35474931_35219951_36449803_36368057_36424128"; __utmz=30149280.1706838836.7.6.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/column/article/612328; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12696; _ga=GA1.2.478668193.1700728840; _ga_Y4GN1R87RG=GS1.1.1706844569.2.1.1706844957.0.0.0; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1709532615%2C%22https%3A%2F%2Ftime.geekbang.org%2Fcolumn%2Farticle%2F612328%22%5D; _pk_ses.100001.8cb4=1; ap_v=0,6.0; __utma=30149280.125988448.1699844939.1709286691.1709532615.11; __utmc=30149280; __utmt=1; dbcl2="126963156:PCmGfH4AjBA"; ck=lbgb; __utmb=30149280.29.5.1709534687162`

var DoubangroupTask = &collect.Task{
	Name:     "find_douban_sun_room",
	WaitTime: 1 * time.Second,
	MaxDepth: 5,
	Cookie:   cookie,
	Rule: collect.RuleTree{
		Root: func() []*collect.Request {
			var roots []*collect.Request
			for i := 0; i < 25; i += 25 {
				url := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
				roots = append(roots, &collect.Request{
					Url:      url,
					Method:   "GET",
					Priority: 1,
					RuleName: "解析网站URL",
				})
			}
			return roots
		},
		Trunk: map[string]*collect.Rule{
			"解析网站URL": {ParseFunc: ParseURL},
			"解析阳台房":   {ParseFunc: GetSumRoom},
		},
	},
}

const urlListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(ctx *collect.Context) collect.ParseResult {
	re := regexp.MustCompile(urlListRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)

	result := collect.ParseResult{}
	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests, &collect.Request{
			Task:     ctx.Req.Task,
			Url:      u,
			Method:   "GET",
			Depth:    ctx.Req.Depth + 1,
			RuleName: "解析阳台房",
		})
	}
	return result
}

const contentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`

func GetSumRoom(ctx *collect.Context) collect.ParseResult {
	re := regexp.MustCompile(contentRe)
	ok := re.Match(ctx.Body)

	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}
	}

	return collect.ParseResult{
		Items: []interface{}{ctx.Req.Url},
	}
}
