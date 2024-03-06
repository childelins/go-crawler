package doubangroup

import (
	"go-crawler/collect"
	"regexp"
)

const urlListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`

func ParseURL(contents []byte, req *collect.Request) collect.ParseResult {
	re := regexp.MustCompile(urlListRe)

	matches := re.FindAllSubmatch(contents, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests, &collect.Request{
			Task:   req.Task,
			Url:    u,
			Method: "GET",
			Depth:  req.Depth + 1,
			ParseFunc: func(b []byte, r *collect.Request) collect.ParseResult {
				return GetContent(b, u)
			},
		})
	}

	return result
}

const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div`

func GetContent(contents []byte, url string) collect.ParseResult {
	re := regexp.MustCompile(ContentRe)
	ok := re.Match(contents)
	if !ok {
		return collect.ParseResult{
			Items: []interface{}{},
		}
	}

	return collect.ParseResult{
		Items: []interface{}{url},
	}
}

// const ContentRe = `<div\s+class="topic-content">(?s:.)*?</div>`

// func GetContent(contents []byte, url string) collect.ParseResult {
// 	re := regexp.MustCompile(ContentRe)
// 	resultStr := re.FindString(string(contents))

// 	r2 := regexp.MustCompile("阳台")

// 	ok := r2.MatchString(resultStr)
// 	if !ok {
// 		return collect.ParseResult{
// 			Items: []interface{}{},
// 		}
// 	}

// 	result := collect.ParseResult{
// 		Items: []interface{}{url},
// 	}

// 	return result
// }
