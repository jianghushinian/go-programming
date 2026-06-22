//go:build ignore
// +build ignore

package main

import (
	"log"
	"os"
)

func main() {
	f, _ := os.Create("demo.log")
	defer f.Close()
	// 日志输出到文件
	logger := log.New(f, "[Debug] - ", log.Lshortfile|log.Lmsgprefix)
	logger.Println("to file")

	logger.SetOutput(os.Stdout)                       // 修改输出位置
	logger.SetPrefix("[Info] - ")                     // 修改前缀
	logger.SetFlags(log.Ldate | log.Ltime | log.LUTC) // 修改日志头
	logger.Println("to stdout with timestamp")
}
