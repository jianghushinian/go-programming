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
	var count int
	// 读写锁，用于保护计数器
	var mu sync.RWMutex

	// 模拟读多写少的场景
	readers := 10                           // 10 个读 goroutine
	writers := 2                            // 2 个写 goroutine
	readInterval := 100 * time.Millisecond  // 读取间隔
	writeInterval := 500 * time.Millisecond // 写入间隔

	// 启动读取协程
	for i := 0; i < readers; i++ {
		go func(id int) {
			for {
				mu.RLock() // 读锁
				val := count
				mu.RUnlock()

				fmt.Printf("reader %d: count = %d\n", id, val)
				time.Sleep(readInterval)
			}
		}(i)
	}

	// 启动写入协程
	for i := 0; i < writers; i++ {
		go func(id int) {
			for j := 0; j < 5; j++ {
				mu.Lock() // 写锁
				count++
				fmt.Printf("writer %d: count = %d\n", id, count)
				mu.Unlock()

				time.Sleep(writeInterval)
			}
		}(i)
	}

	// 等待写入完成
	time.Sleep(3 * time.Second)
	fmt.Printf("count: %d\n", count)
}
