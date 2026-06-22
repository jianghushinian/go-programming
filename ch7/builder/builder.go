package builder

import "time"

// Server 要被构造的复杂对象
type Server struct {
	Addr    string
	Port    int
	Timeout time.Duration
}

// Builder 建造者抽象接口
type Builder interface {
	SetAddr(addr string)
	SetPort(port int)
	SetTimeout(d time.Duration)
	Build() *Server
}

// ServerBuilder 具体建造者
type ServerBuilder struct {
	s *Server
}

func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{s: &Server{}}
}

func (b *ServerBuilder) SetAddr(addr string)        { b.s.Addr = addr }
func (b *ServerBuilder) SetPort(port int)           { b.s.Port = port }
func (b *ServerBuilder) SetTimeout(d time.Duration) { b.s.Timeout = d }
func (b *ServerBuilder) Build() *Server             { return b.s }

// Director 指挥者
type Director struct {
	builder Builder
}

func NewDirector(b Builder) *Director {
	return &Director{builder: b}
}

// ConstructServer 指挥构造 Server 对象，内部可以控制构造顺序
func (d *Director) ConstructServer() *Server {
	d.builder.SetAddr("127.0.0.1")
	d.builder.SetPort(8080)
	d.builder.SetTimeout(10 * time.Second)
	return d.builder.Build()
}
