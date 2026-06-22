package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	genericapiserver "k8s.io/apiserver/pkg/server"
)

func main() {
	// 初始化 HTTP 服务
	srv := &http.Server{
		Addr: ":8000",
	}

	// 注册路由处理函数
	http.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		duration, err := time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		time.Sleep(duration)
		_, _ = w.Write([]byte("Welcome HTTP Server"))
	})

	// 开启新的 goroutine 启动 HTTP 服务，避免阻塞主 goroutine
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		log.Println("Stopped serving new connections")
	}()

	// 等待终止信号并优雅退出
	waitForShutdown(srv)
}

// waitForShutdown 等待终止信号并优雅退出服务器
func waitForShutdown(srv *http.Server) {
	// 等待终止信号
	<-genericapiserver.SetupSignalHandler()
	log.Println("Shutdown Server...")

	// 创建带超时的上下文用于优雅退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 执行优雅退出
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
	log.Println("HTTP server graceful shutdown completed")
}
