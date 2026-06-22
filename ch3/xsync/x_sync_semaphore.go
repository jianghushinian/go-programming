//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/semaphore"
)

// runWorkerPool 用信号量控制最大并发的工作池
func runWorkerPool(maxWorkers int, out []int) error {
	// 创建信号量，限制最大并发数
	sem := semaphore.NewWeighted(int64(maxWorkers))
	ctx := context.TODO()

	// 启动任务
	for i := range out {
		// 获取信号量
		if err := sem.Acquire(ctx, 1); err != nil {
			return fmt.Errorf("failed to acquire semaphore: %w", err)
		}

		// 启动新的 goroutine 执行任务
		go func(i int) {
			defer sem.Release(1) // 任务完成后释放信号量
			out[i] = fn(i + 1)   // 模拟任务计算
		}(i)
	}

	// 等待所有任务完成：获取所有信号量，确保所有 goroutine 完成
	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		return fmt.Errorf("failed to acquire semaphore for completion: %w", err)
	}
	return nil
}

// fn 模拟一个耗时任务，返回结果
func fn(n int) int {
	time.Sleep(1 * time.Second) // 模拟耗时
	fmt.Printf("run %d at %v\n", n, time.Now().Format("15:04:05"))
	return n * 2 // 返回任务执行结果
}

func main() {
	maxWorkers := 2       // 最大并发数
	out := make([]int, 6) // 总任务数量
	// 运行工作池
	if err := runWorkerPool(maxWorkers, out); err != nil {
		log.Fatalf("Error running worker pool: %v", err)
	}
	fmt.Println(out) // 输出计数结果
}
