package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// 一个任务实例
type Task struct {
	Url         string
	Cookie      string
	WaitTime    time.Duration
	Reload      bool // 网站是否可以重复爬取
	MaxDepth    int  // 爬取最大深度
	Visited     map[string]bool
	VisitedLock sync.Mutex
	RootReq     *Request // 初始请求
	Fetcher     Fetcher
}

// 单个请求
type Request struct {
	unique    string
	Task      *Task
	Url       string
	Method    string
	Depth     int // 当前深度
	Priority  int // 优先级队列
	ParseFunc func([]byte, *Request) ParseResult
}

// 爬取结果
type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("max depth limit reached")
	}

	return nil
}

// 请求的唯一识别码
func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}
