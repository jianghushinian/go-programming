//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"
)

type Server struct {
	Addr    string
	Timeout time.Duration
}

type ServerOptions struct {
	Addr    string
	Timeout time.Duration
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		Addr:    "127.0.0.1",
		Timeout: 5 * time.Second,
	}
}

func NewServerWithOptions(opts *ServerOptions) *Server {
	return &Server{
		Addr:    opts.Addr,
		Timeout: opts.Timeout,
	}
}

func main() {
	// 使用默认配置
	s1 := NewServerWithOptions(NewServerOptions())

	// 自定义配置
	opts := NewServerOptions()
	opts.Timeout = 10 * time.Second
	s2 := NewServerWithOptions(opts)

	fmt.Println(s1)
	fmt.Println(s2)
}
