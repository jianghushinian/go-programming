package main

import (
	"testing"

	"golang.org/x/sync/singleflight"
)

func TestDoChanPanic(t *testing.T) {
	g := new(singleflight.Group)

	// 测试函数
	fn := func() (interface{}, error) {
		// 此 panic 无法被 recover
		panic("test panic")
	}

	ch := g.DoChan("key", fn)

	// 等待结果
	result := <-ch

	// 验证错误类型
	if result.Err != nil {
		t.Fatalf("预期 PanicError，但收到 %T", result.Err)
	}

	t.Log("-----------")
}
