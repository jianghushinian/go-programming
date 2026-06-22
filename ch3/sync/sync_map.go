//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
)

func main() {
	var s sync.Map

	// 写操作
	s.Store("name", "江湖十年")
	s.Store("addr", "杭州")

	// 读操作
	if value, ok := s.Load("name"); ok {
		fmt.Println("name:", value)
	}

	// 删除操作
	s.Delete("name")

	// 遍历操作
	s.Range(func(key, value any) bool {
		fmt.Printf("%s: %s\n", key, value)
		return true // 继续遍历
	})
}
