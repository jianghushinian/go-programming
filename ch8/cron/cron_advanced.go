//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type MyJob struct {
	name  string
	count int
}

func (j *MyJob) Run() {
	j.count++
	if j.count == 2 {
		panic("xxxx")
	}
	fmt.Printf("任务 %s，第 %d 次执行\n", j.name, j.count)
}

type cronLogger struct{}

func (l *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	// 自定义输出，比如用 slog、zap 等
	log.Println("INFO:", msg, keysAndValues)
}

func (l *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Println("ERROR:", msg, keysAndValues, "err:", err)
}

func main() {
	logger := &cronLogger{}
	c := cron.New(
		cron.WithSeconds(),
		cron.WithLogger(logger),
		cron.WithChain(cron.SkipIfStillRunning(logger), cron.Recover(logger)),
		// Bug: https://github.com/robfig/cron/commit/bc59245fe10efaed9d51b56900192527ed733435
		// cron.WithChain(cron.Recover(logger), cron.SkipIfStillRunning(logger)),
	)
	id, err := c.AddJob("@every 5s", &MyJob{name: "江湖十年"})
	if err != nil {
		log.Fatalf("添加任务失败: %v", err)
	}
	c.Start()
	defer c.Stop()

	time.Sleep(30 * time.Second)
	c.Remove(id)
	time.Sleep(5 * time.Second)
}
