//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group
	// 限制最大并发任务数为 2
	g.SetLimit(2)

	// 并发运行 3 个任务
	for i := 1; i <= 3; i++ {
		if g.TryGo(func() error {
			fmt.Printf("Goroutine %d is starting\n", i)
			time.Sleep(2 * time.Second) // 模拟工作耗时
			fmt.Printf("Goroutine %d is done\n", i)
			return nil
		}) { // 运行成功
			fmt.Printf("Goroutine %d started successfully\n", i)
		} else { // 达到并发限制，打印提示
			fmt.Printf("Goroutine %d could not start (limit reached)\n", i)
		}
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Encountered an error: %v\n", err)
	}
	fmt.Println("All goroutines complete.")
}
