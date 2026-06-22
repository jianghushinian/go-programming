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

func NewServer(addr string) *Server {
	// 提供默认的超时时间
	return &Server{Addr: addr, Timeout: 5 * time.Second}
}

func NewServerWithTimeout(addr string, timeout time.Duration) *Server {
	return &Server{Addr: addr, Timeout: timeout}
}

func main() {
	server := NewServer("127.0.0.1:8080")
	fmt.Println(server)

	server = NewServerWithTimeout("127.0.0.1:8080", time.Second)
	fmt.Println(server)
}
