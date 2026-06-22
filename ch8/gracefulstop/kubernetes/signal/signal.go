package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var onlyOneSignalHandler = make(chan struct{})
var shutdownHandler chan os.Signal
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func SetupSignalHandler() <-chan struct{} {
	return SetupSignalContext().Done()
}

func SetupSignalContext() context.Context {
	close(onlyOneSignalHandler)

	shutdownHandler = make(chan os.Signal, 2) // 缓冲区长度为 2

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(shutdownHandler, shutdownSignals...)
	go func() {
		<-shutdownHandler // 第一次信号：触发取消
		cancel()
		<-shutdownHandler // 第二次信号：强制退出
		os.Exit(1)
	}()

	return ctx
}
