package option

import "time"

type Server struct {
	Addr    string
	Timeout time.Duration
}

// Option 定义选项类型
type Option func(*Server)

// NewServer 构造函数，应用所有配置选项
func NewServer(opts ...Option) *Server {
	s := &Server{
		Addr:    "127.0.0.1",     // 默认值
		Timeout: 5 * time.Second, // 默认值
	}
	for _, opt := range opts {
		opt(s) // 应用每个选项，替代默认值
	}
	return s
}

// WithAddr 设置地址
func WithAddr(addr string) Option {
	return func(s *Server) {
		s.Addr = addr
	}
}

// WithTimeout 设置超时时间
func WithTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.Timeout = d
	}
}
