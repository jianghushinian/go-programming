package builder

import "time"

type SimpleServerBuilder struct {
	s *Server
}

func NewSimpleServerBuilder() *SimpleServerBuilder {
	return &SimpleServerBuilder{
		s: &Server{
			Addr:    "127.0.0.1",
			Port:    8080,
			Timeout: 10 * time.Second,
		},
	}
}

func (b *SimpleServerBuilder) Addr(a string) *SimpleServerBuilder {
	b.s.Addr = a
	return b
}

func (b *SimpleServerBuilder) Port(p int) *SimpleServerBuilder {
	b.s.Port = p
	return b
}

func (b *SimpleServerBuilder) Timeout(d time.Duration) *SimpleServerBuilder {
	b.s.Timeout = d
	return b
}

func (b *SimpleServerBuilder) Build() *Server {
	return b.s
}
