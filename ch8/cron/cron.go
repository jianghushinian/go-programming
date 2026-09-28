package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	// 创建一个新的 Cron 实例
	c := cron.New(cron.WithSeconds())

	// 添加定时任务
	c.AddFunc("* * * * * *", func() { job(1) })
	c.AddFunc("*/5 * * * * *", func() { job(5) })

	c.Start()      // 启动计划任务
	defer c.Stop() // 关闭计划任务，这不会关闭已经在执行中的任务

	select {} // 阻塞语句，保持程序运行
}

func job(duration int) {
	fmt.Printf("每 %d 秒执行的任务: %s\n", duration, time.Now().Format("2006-01-02 15:04:05"))
}
