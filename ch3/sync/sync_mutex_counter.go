//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 计数器初始值
	var count = 0
	// 互斥锁，用于保护计数器
	var mu sync.Mutex

	// 开启 2 个 goroutine 并发修改 count
	for i := 0; i < 2; i++ {
		go func() {
			for j := 0; j < 100000; j++ {
				mu.Lock() // 加锁
				count++
				mu.Unlock() // 解锁
			}
		}()
	}

	// 等待 2 个 goroutine 执行完成
	time.Sleep(1 * time.Second)

	// 输出最终计数器的值
	fmt.Println(count)
}
