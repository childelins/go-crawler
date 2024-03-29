package main

import (
	"fmt"
	"sort"
)

// 计算机课程和其前序课程的映射关系
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

// 深度优先搜索
func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)

	// visitAll 会使用递归计算最前序的课程，并添加到列表的 order 中，这就保证了课程的先后顺序
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				// fmt.Println("seen:", seen)
				seen[item] = true

				// 深度递归
				visitAll(m[item])

				// 跳出递归后，开始添加到 order
				order = append(order, item)
				// fmt.Println("order:", order)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	// fmt.Println("keys:", keys)
	visitAll(keys)
	return order
}
