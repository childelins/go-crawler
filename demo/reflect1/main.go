package main

import "fmt"

type Student struct {
	Age  int
	Name string
}

func (s *Student) CreateSQL() string {
	sql := fmt.Sprintf("insert into student values(%d, %s)", s.Age, s.Name)
	return sql
}

func main() {
	s := Student{
		Age:  20,
		Name: "jonson",
	}

	fmt.Println(s.CreateSQL())
}
