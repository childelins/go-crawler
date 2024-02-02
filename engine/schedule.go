package engine

import (
	"go-crawler/collect"

	"go.uber.org/zap"
)

type ScheduleEngine struct {
	requestCh chan *collect.Request
	workerCh  chan *collect.Request
	WorkCount int
	Fetcher   collect.Fetcher
	Logger    *zap.Logger
	out       chan collect.ParseResult
	Seeds     []*collect.Request
}

func (s *ScheduleEngine) Run() {
	requesCh := make(chan *collect.Request)
	workerCh := make(chan *collect.Request)
	out := make(chan collect.ParseResult)

	s.requestCh = requesCh
	s.workerCh = workerCh
	s.out = out

	go s.Schedule()
	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}

	s.HandleResult()
}

func (s *ScheduleEngine) Schedule() {
	reqQueue := s.Seeds
	go func() {
		for {
			var req *collect.Request
			var ch chan *collect.Request

			if len(reqQueue) > 0 {
				req = reqQueue[0]
				reqQueue = reqQueue[1:]
				ch = s.workerCh
			}

			select {
			case r := <-s.requestCh:
				// 接收请求，存放在 reqQueue 队列中
				reqQueue = append(reqQueue, r)
			case ch <- req:
				// 从 reqQueue 取一个任务，分发到 worker channel，等待 worker 处理
			}
		}
	}()
}

func (s *ScheduleEngine) CreateWork() {
	for {
		r := <-s.workerCh
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.Logger.Error("can't fetch ", zap.Error(err))
			continue
		}

		result := r.ParseFunc(body, r)
		s.out <- result
	}
}

func (s *ScheduleEngine) HandleResult() {
	for {
		select {
		case result := <-s.out:
			for _, req := range result.Requests {
				s.requestCh <- req
			}
			for _, item := range result.Items {
				// todo: store
				s.Logger.Sugar().Info("get result", item)
			}
		}
	}

}
