package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("enter")
	quit := make(chan os.Signal, 1)
	// 注册需要关注的信号：SIGINT、SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞当前 goroutine 等待信号
	sig := <-quit
	fmt.Printf("\nreceived signal: %d-%s\n", sig, sig)
	fmt.Println("exit")
}
