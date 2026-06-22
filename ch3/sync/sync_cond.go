//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	cond := sync.NewCond(&mu) // 创建 cond，需要传入一个实现 Locker 接口的锁
	var ready bool            // 共享条件变量

	// 启动一个 goroutine 等待条件
	go func() {
		fmt.Println("goroutine: 等待条件满足...")
		cond.L.Lock() // 1. 进入临界区，加锁
		for !ready {  // 2. 循环检查条件 (防止虚假唤醒)
			cond.Wait() // 3. 等待时会暂时释放锁，被唤醒后重新获取锁
		}
		fmt.Println("goroutine: 条件已满足!")
		cond.L.Unlock() // 4. 操作完成，解锁
	}()

	time.Sleep(2 * time.Second) // 模拟主 goroutine 做一些准备工作

	cond.L.Lock()
	ready = true // 5. 更改条件变量
	fmt.Println("main: 条件已设置，通知等待者")
	cond.Signal() // 6. 发送信号，唤醒一个等待的 goroutine
	cond.L.Unlock()

	time.Sleep(time.Second) // 等待一下，让子 goroutine 完成输出
}
