package main

import (
	"embed"
	"io/fs"
	"net/http"
	"sync"
)

//go:embed static
var staticFS embed.FS

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// 服务 1：直接暴露 static 目录
	// 使用 go:embed 实现静态文件服务
	go func() {
		defer wg.Done()
		_ = http.ListenAndServe(":8001", http.FileServer(http.FS(staticFS)))
	}()

	// 服务 2：通过 fs.Sub 去掉 URL 中的 /static 前缀
	go func() {
		defer wg.Done()
		fsSub, _ := fs.Sub(staticFS, "static")
		_ = http.ListenAndServe(":8002", http.FileServer(http.FS(fsSub)))
	}()

	wg.Wait()
}
