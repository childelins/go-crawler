package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go-crawler/collector"
	"time"
)

type Property struct {
	Name     string        `json:"name"` // 任务名称，应保证唯一性
	Url      string        `json:"url"`
	Cookie   string        `json:"cookie"`
	WaitTime time.Duration `json:"wait_time"`
	Reload   bool          `json:"reload"`    // 网站是否可以重复爬取
	MaxDepth int64         `json:"max_depth"` // 爬取最大深度
}

// 一个任务实例
type Task struct {
	Property

	// Visited     map[string]bool
	// VisitedLock sync.Mutex
	// Fetcher Fetcher
	Rule RuleTree // 规则条件
}

type Context struct {
	Body []byte
	Req  *Request
}

func (c *Context) GetRule(ruleName string) *Rule {
	return c.Req.Task.Rule.Trunk[ruleName]
}

func (c *Context) Output(data interface{}) *collector.OutputData {
	res := &collector.OutputData{}
	res.Data = make(map[string]interface{})
	res.Data["Rule"] = c.Req.RuleName
	res.Data["Data"] = data
	res.Data["Url"] = c.Req.Url
	res.Data["Time"] = time.Now().Format("2006-01-02 15:04:05")
	return res
}

// 单个请求
type Request struct {
	// unique    string
	Task     *Task // 指向具体任务
	Url      string
	Method   string
	Depth    int64  // 当前深度
	Priority int64  // 优先级队列
	RuleName string // 规则名称
	TmpData  *Temp  // 临时缓存数据
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
