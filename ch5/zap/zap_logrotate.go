//go:build ignore
// +build ignore

package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"          // Zap 核心组件
	"gopkg.in/natefinch/lumberjack.v2" // 日志轮转库
)

func main() {
	// 将 io.Writer 转换为支持 Sync 方法的 WriteSyncer
	file := zapcore.AddSync(&lumberjack.Logger{ // 带轮转功能的日志写入器
		Filename:   "/var/log/app.log",
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
	})

	// JSON 格式的日志编码器
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// 日志核心对象
	core := zapcore.NewCore(encoder, file, zap.InfoLevel)
	// 创建 Logger 实例
	logger := zap.New(core)
	defer logger.Sync() // 确保程序退出前刷新缓冲区

	logger.Info("Server started")

	// 创建控制台输出目标对象（带线程安全锁）
	console := zapcore.Lock(os.Stdout)
	// 创建多路复用核心（Tee）
	tee := zapcore.NewTee(
		zapcore.NewCore(encoder, console, zap.DebugLevel),
		zapcore.NewCore(encoder, file, zap.InfoLevel),
	)
	// 创建 Logger 实例
	logger = zap.New(tee)
	defer logger.Sync() // 确保程序退出前刷新缓冲区

	logger.Debug("logging for debug")
	logger.Info("logging for info")
}
