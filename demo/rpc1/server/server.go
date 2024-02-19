package server

import (
	"errors"
)

// 请求参数
type Args struct {
	A, B int
}

// 响应参数
type Quotient struct {
	Quo, Rem int
}

// 服务对象
type Arith int

func (t *Arith) Multiply(arg *Args, reply *int) error {
	*reply = arg.A * arg.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}

	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}
