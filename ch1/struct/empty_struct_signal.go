//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})

	go func() {
		time.Sleep(1 * time.Second) // 模拟执行业务逻辑
		fmt.Printf("goroutine done\n")
		done <- struct{}{} // 发送完成信号
	}()

	fmt.Printf("waiting...\n")
	<-done // 等待子协程完成
	fmt.Printf("main exit\n")
}
