package engine

import (
	"go-crawler/collect"
	"go-crawler/parse/doubangroup"
	"sync"

	"go.uber.org/zap"
)

func init() {
	Store.Add(doubangroup.DoubangroupTask)
}

// 全局蜘蛛种类实例
var Store = &CrawlerStore{
	list: []*collect.Task{},
	hash: map[string]*collect.Task{},
}

type CrawlerStore struct {
	list []*collect.Task
	hash map[string]*collect.Task
}

func (c *CrawlerStore) Add(task *collect.Task) {
	c.hash[task.Name] = task
	c.list = append(c.list, task)
}

type Crawler struct {
	out         chan collect.ParseResult // 爬取结果
	Visited     map[string]bool          // 过滤重复请求
	VisitedLock sync.Mutex

	failures    map[string]*collect.Request // 失败请求id -> 失败请求
	failureLock sync.Mutex

	options
}

type Scheduler interface {
	Schedule()
	Push(...*collect.Request)
	Pull() *collect.Request
}

type Schedule struct {
	requestCh   chan *collect.Request
	workerCh    chan *collect.Request
	priReqQueue []*collect.Request // 优先队列
	reqQueue    []*collect.Request // 普通队列
	// Logger      *zap.Logger
}

func NewEngine(opts ...Option) *Crawler {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}

	e := &Crawler{}
	e.Visited = make(map[string]bool, 100)
	e.out = make(chan collect.ParseResult)
	e.failures = make(map[string]*collect.Request)
	e.options = options
	return e
}

func NewSchedule() *Schedule {
	s := &Schedule{}
	requestCh := make(chan *collect.Request)
	workerCh := make(chan *collect.Request)
	s.requestCh = requestCh
	s.workerCh = workerCh
	return s
}

func (e *Crawler) Run() {
	go e.Schedule()
	for i := 0; i < e.WorkCount; i++ {
		go e.CreateWork()
	}
	e.HandleResult()
}

func (s *Schedule) Push(reqs ...*collect.Request) {
	for _, req := range reqs {
		s.requestCh <- req
	}
}

func (s *Schedule) Pull() *collect.Request {
	r := <-s.workerCh
	return r
}

func (s *Schedule) Schedule() {
	var req *collect.Request
	var ch chan *collect.Request
	for {
		// 先处理优先队列
		if req == nil && len(s.priReqQueue) > 0 {
			req = s.priReqQueue[0]
			s.priReqQueue = s.priReqQueue[1:]
			ch = s.workerCh
		}
		// 再处理普通队列
		if req == nil && len(s.reqQueue) > 0 {
			req = s.reqQueue[0]
			s.reqQueue = s.reqQueue[1:]
			ch = s.workerCh
		}

		select {
		case r := <-s.requestCh:
			// 根据请求的优先级，进对应队列
			if r.Priority > 0 {
				s.priReqQueue = append(s.priReqQueue, r)
			} else {
				s.reqQueue = append(s.reqQueue, r)
			}
		case ch <- req:
			// req有值，此时 ch = workerCh， 即往 workerCh 写入数据
			// 重置 req，进入下一轮 for 循环
			req = nil
			ch = nil
		}
	}
}

func (e *Crawler) Schedule() {
	var reqs []*collect.Request
	for _, seed := range e.Seeds {
		// must store value
		task := Store.hash[seed.Name]
		// 获取初始化任务
		rootreqs := task.Rule.Root()
		for _, req := range rootreqs {
			// 此时赋值 task
			req.Task = task
		}

		reqs = append(reqs, rootreqs...)
	}

	go e.scheduler.Schedule()
	go e.scheduler.Push(reqs...)
}

func (e *Crawler) CreateWork() {
	for {
		req := e.scheduler.Pull()
		if err := req.Check(); err != nil {
			e.Logger.Error("check failed", zap.Error(err))
			continue
		}

		if !req.Task.Reload && e.HasVisited(req) {
			e.Logger.Debug("request has visited", zap.String("url:", req.Url))
		}
		e.StoreVisited(req)

		body, err := e.Fetcher.Get(req)
		if err != nil {
			e.Logger.Error("can't fetch ", zap.Error(err), zap.String("url", req.Url))
			e.SetFailure(req)
			continue
		}
		if len(body) < 6000 {
			e.Logger.Error("can't fetch ", zap.Int("length", len(body)), zap.String("url", req.Url))
			e.SetFailure(req)
			continue
		}

		// 获取当前任务对应的规则
		rule := req.Task.Rule.Trunk[req.RuleName]
		// 内容解析
		result := rule.ParseFunc(&collect.Context{
			Body: body,
			Req:  req,
		})

		if len(result.Requests) > 0 {
			// 新的任务加入队列中
			go e.scheduler.Push(result.Requests...)
		}

		e.out <- result
	}
}

func (e *Crawler) HandleResult() {
	for {
		select {
		case result := <-e.out:
			for _, item := range result.Items {
				// todo: store
				e.Logger.Sugar().Info("get result: ", item)
			}
		}
	}
}

func (e *Crawler) HasVisited(r *collect.Request) bool {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()
	unique := r.Unique()
	return e.Visited[unique]
}

func (e *Crawler) StoreVisited(reqs ...*collect.Request) {
	e.VisitedLock.Lock()
	defer e.VisitedLock.Unlock()

	for _, r := range reqs {
		unique := r.Unique()
		e.Visited[unique] = true
	}
}

func (e *Crawler) SetFailure(req *collect.Request) {
	if !req.Task.Reload {
		e.VisitedLock.Lock()
		delete(e.Visited, req.Unique())
		e.VisitedLock.Unlock()
	}

	e.failureLock.Lock()
	defer e.failureLock.Unlock()
	if _, ok := e.failures[req.Unique()]; !ok {
		// 首次失败时，再重新执行一次
		e.failures[req.Unique()] = req
		e.scheduler.Push(req)
	}

	// todo: 失败2次，加载到失败队列中
}
