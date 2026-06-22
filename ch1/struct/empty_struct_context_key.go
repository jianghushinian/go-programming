package main

import (
	"context"
	"fmt"
)

type requestIdKey struct{}

func main() {
	// 正确示例
	{
		ctx := context.Background()

		// 设置值
		ctx = context.WithValue(ctx, requestIdKey{}, "req-123")
		// 获取值
		fmt.Printf("request-id: %s\n", ctx.Value(requestIdKey{}))
	}

	// 错误示例
	{
		ctx := context.Background()

		key1 := struct{}{}
		ctx = context.WithValue(ctx, key1, "data1")
		fmt.Printf("key1 data: %s\n", ctx.Value(key1))

		key2 := struct{}{}
		ctx = context.WithValue(ctx, key2, "data2")
		fmt.Printf("key2 data: %s\n", ctx.Value(key2))

		// 再次查看 key1 对应的 value
		fmt.Printf("key1 data: %s\n", ctx.Value(key1))
	}
}
