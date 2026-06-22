//go:build ignore
// +build ignore

package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger) // 将标准库 log 重定向到 zap
	log.Println("this goes to zap")
}
