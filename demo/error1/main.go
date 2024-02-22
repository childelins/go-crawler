package main

import (
	"errors"
	"fmt"
)

var ErrDivideZero = errors.New("divide by zero")

func main() {
	val, err := Cal(10, 0)
	if err != nil {
		if ok := errors.Is(err, ErrDivideZero); ok {
			fmt.Printf("divide err: %v\n", err)
		} else {
			fmt.Printf("other err: %v\n", err)
		}
		return
	}

	fmt.Println("val:", val)
}

func Cal(a, b int) (int, error) {
	val, err := Divide(a, b)
	if err != nil {
		// %w 实现了wrap error，但不具备打印堆栈信息功能
		return 0, fmt.Errorf("cal error: %w", err)
	}

	return val, nil
}

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideZero
	}

	return a / b, nil
}
