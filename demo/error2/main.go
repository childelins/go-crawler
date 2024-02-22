package main

import (
	"fmt"

	"errors"

	xerrors "github.com/pkg/errors"
)

var ErrDivideZero = errors.New("divide by zero")

func main() {
	val, err := Cal(10, 0)
	if err != nil {
		if ok := errors.Is(err, ErrDivideZero); ok {
			// fmt.Println("original error:", errors.Unwrap(err))
			fmt.Printf("divide err: %+v\n", err)
		} else {
			fmt.Printf("other err: %+v\n", err)
		}
		return
	}

	fmt.Println("val:", val)
}

func Cal(a, b int) (int, error) {
	val, err := Divide(a, b)
	if err != nil {
		// 变化点: 实现了wrap error，且具备打印堆栈信息功能
		return 0, xerrors.Wrap(err, "cal error")
		// return 0, xerrors.WithMessage(err, "cal error")
	}

	return val, nil
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideZero
	}

	return a / b, nil
}
