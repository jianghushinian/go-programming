//go:build ignore
// +build ignore

package main

import "go.uber.org/zap"

func main() {
	atom := zap.NewAtomicLevelAt(zap.InfoLevel) // 支持动态调整级别
	cfg := zap.Config{
		Level:            atom,
		Development:      false,
		Sampling:         &zap.SamplingConfig{Initial: 1, Thereafter: 2},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr", "/var/log/app.json"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := cfg.Build() // 构建 Logger 对象
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	atom.SetLevel(zap.DebugLevel) // 在运行时修改日志级别
	for i := 0; i < 3; i++ {
		logger.Debug("logging for debug")
	}
}
