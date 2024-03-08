package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"
)

// 一个任务实例
type Task struct {
	Name     string // 用户界面显示的名称（应保证唯一性）
	Url      string
	Cookie   string
	WaitTime time.Duration
	Reload   bool // 网站是否可以重复爬取
	MaxDepth int  // 爬取最大深度
	// Visited     map[string]bool
	// VisitedLock sync.Mutex
	// Fetcher Fetcher
	Rule RuleTree // 规则条件
}

type Context struct {
	Body []byte
	Req  *Request
}

// 单个请求
type Request struct {
	// unique    string
	Task     *Task // 指向具体任务
	Url      string
	Method   string
	Depth    int    // 当前深度
	Priority int    // 优先级队列
	RuleName string // 规则名称
}

// 爬取结果
type ParseResult struct {
	Requests []*Request    // 下一批请求
	Items    []interface{} // 当前请求爬取结果
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
