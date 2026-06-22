//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu1, mu2 sync.Mutex

	// goroutine 1
	go func() {
		mu1.Lock()
		fmt.Println("goroutine 1 locked")

		time.Sleep(100 * time.Millisecond) // 模拟工作

		mu2.Lock() // 尝试获取锁 2（会阻塞）
		fmt.Println("goroutine 1 lock 2")

		mu2.Unlock()
		mu1.Unlock()
	}()

	// goroutine 2
	go func() {
		mu2.Lock()
		fmt.Println("goroutine 2 locked")

		time.Sleep(100 * time.Millisecond) // 模拟工作

		mu1.Lock() // 尝试获取锁 1（会阻塞）
		fmt.Println("goroutine 2 lock 1")

		mu1.Unlock()
		mu2.Unlock()
	}()

	select {}
}
