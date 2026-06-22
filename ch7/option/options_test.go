package option

import (
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	// 使用默认配置
	s1 := NewServer()
	t.Log(s1)

	// 仅修改 Addr
	s2 := NewServer(WithAddr("localhost"))
	t.Log(s2)

	// 修改多个配置
	s3 := NewServer(
		WithAddr("0.0.0.0"),
		WithTimeout(10*time.Second),
	)
	t.Log(s3)
}
