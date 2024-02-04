package collect

import (
	"errors"
	"time"
)

type Request struct {
	Url       string
	Cookie    string
	WaitTime  time.Duration
	Depth     int
	MaxDepth  int
	ParseFunc func([]byte, *Request) ParseResult
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

func (r *Request) Check() error {
	if r.Depth > r.MaxDepth {
		return errors.New("max depth limit reached")
	}

	return nil
}
