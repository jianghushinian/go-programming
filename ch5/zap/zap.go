//go:build ignore
// +build ignore

package main

import (
	"go.uber.org/zap"
)

func main() {
	// 生产配置（JSON、INFO 级别）
	logger, _ := zap.NewProduction()
	defer logger.Sync() // 程序退出前将缓冲区中日志写入目标位置
	logger.Info("production",
		zap.String("url", "https://jianghushinian.cn/"),
		zap.Int("attempt", 3))

	// 开发配置（人类可读）
	devLogger, _ := zap.NewDevelopment()
	defer devLogger.Sync()
	devLogger.Debug("dev", zap.String("mode", "dev"))

	// 更方便友好的 API
	sugar := logger.Sugar()
	defer sugar.Sync()
	sugar.Infow("user login",
		"user", 42,
		"ip", "10.0.0.1",
	)
}
