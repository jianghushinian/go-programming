//go:build ignore
// +build ignore

package main

import (
	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
)

func LogRotateHook() logrus.Hook {
	hook, _ := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			Filename:   "app.log", // 统一的日志文件
			MaxSize:    100,       // 100MB 后轮转
			MaxBackups: 10,        // 保留 10 个备份
			MaxAge:     30,        // 保留 30 天
			Compress:   true,      // 压缩旧日志
		},
		logrus.InfoLevel,        // 从 Info 级别开始记录
		&logrus.JSONFormatter{}, // 使用 JSON 格式
		nil,
	)
	return hook
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.AddHook(LogRotateHook())

	logrus.Info("Server started")
}
