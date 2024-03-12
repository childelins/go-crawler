package doubanbook

import (
	"go-crawler/collect"
	"regexp"
	"strconv"
	"time"
)

const cookie = `bid=qk-KbS-ffCg; douban-fav-remind=1; _pk_id.100001.8cb4=2b9cc6a45e496b68.1699844938.; ll="118201"; viewed="1007305_35196328_35474931_35219951_36449803_36368057_36424128"; __utmz=30149280.1706838836.7.6.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/column/article/612328; push_noty_num=0; push_doumail_num=0; __utmv=30149280.12696; _ga=GA1.2.478668193.1700728840; _ga_Y4GN1R87RG=GS1.1.1706844569.2.1.1706844957.0.0.0; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1709532615%2C%22https%3A%2F%2Ftime.geekbang.org%2Fcolumn%2Farticle%2F612328%22%5D; _pk_ses.100001.8cb4=1; ap_v=0,6.0; __utma=30149280.125988448.1699844939.1709286691.1709532615.11; __utmc=30149280; __utmt=1; dbcl2="126963156:PCmGfH4AjBA"; ck=lbgb; __utmb=30149280.29.5.1709534687162`

var DoubanBookTask = &collect.Task{
	Property: collect.Property{
		Name:     "douban_book_list",
		WaitTime: 1 * time.Second,
		MaxDepth: 5,
		Cookie:   cookie,
	},
	Rule: collect.RuleTree{
		Root: func() ([]*collect.Request, error) {
			roots := []*collect.Request{
				{
					Url:      "https://book.douban.com",
					Method:   "GET",
					Priority: 1,
					RuleName: "数据tag",
				},
			}
			return roots, nil
		},
		Trunk: map[string]*collect.Rule{
			"数据tag": {ParseFunc: ParseTag},
			"书籍列表":  {ParseFunc: ParseBookList},
			"书籍简介":  {ParseFunc: ParseBookDetail},
		},
	},
}

// <a href="/tag/小说" class="tag">小说</a>
const regexpStr = `<a href="([^"]+)" class="tag">([^<]+)</a>`

// 解析豆瓣读书 tag 标签
func ParseTag(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(regexpStr)
	matches := re.FindAllSubmatch(ctx.Body, -1)

	result := collect.ParseResult{}
	for _, m := range matches {
		result.Requests = append(result.Requests, &collect.Request{
			Task:     ctx.Req.Task,
			Url:      "https://book.douban.com" + string(m[1]),
			Method:   "GET",
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍列表",
		})
	}

	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	result.Requests = result.Requests[:1]
	return result, nil
}

const BooklistRe = `<a.*?href="([^"]+)" title="([^"]+)"`

// 解析书籍列表
func ParseBookList(ctx *collect.Context) (collect.ParseResult, error) {
	re := regexp.MustCompile(BooklistRe)
	matches := re.FindAllSubmatch(ctx.Body, -1)

	result := collect.ParseResult{}
	for _, m := range matches {
		req := &collect.Request{
			Task:     ctx.Req.Task,
			Url:      string(m[1]),
			Method:   "GET",
			Depth:    ctx.Req.Depth + 1,
			RuleName: "书籍简介",
		}

		req.TmpData = &collect.Temp{}
		req.TmpData.Set("book_name", string(m[2]))
		result.Requests = append(result.Requests, req)
	}

	// 在添加limit之前，临时减少抓取数量,防止被服务器封禁
	result.Requests = result.Requests[:3]
	return result, nil
}

var autoRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>([^<]+)</a>`)
var publicRe = regexp.MustCompile(`<span class="pl">出版社:</span>([^<]+)<br/>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> ([^<]+)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span>([^<]+)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average">([^<]+)</strong>`)
var intoRe = regexp.MustCompile(`<div class="intro">[\d\D]*?<p>([^<]+)</p></div>`)

func ParseBookDetail(ctx *collect.Context) (collect.ParseResult, error) {
	bookName := ctx.Req.TmpData.Get("book_name")
	page, _ := strconv.Atoi(ExtraString(ctx.Body, pageRe))

	book := map[string]interface{}{
		"书名":  bookName,
		"作者":  ExtraString(ctx.Body, autoRe),
		"页数":  page,
		"出版社": ExtraString(ctx.Body, publicRe),
		"得分":  ExtraString(ctx.Body, scoreRe),
		"价格":  ExtraString(ctx.Body, priceRe),
		"简介":  ExtraString(ctx.Body, intoRe),
	}

	data := ctx.Output(book)

	result := collect.ParseResult{
		Items: []interface{}{data},
	}

	return result, nil
}

func ExtraString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
