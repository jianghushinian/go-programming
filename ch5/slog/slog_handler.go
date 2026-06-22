//go:build ignore
// +build ignore

package main

import (
	"context"
	"log/slog"
	"os"
)

// ContextHandler 包装现有的 Handler，增加上下文属性提取功能
type ContextHandler struct {
	slog.Handler
}

// Handle 方法重写，在调用底层 Handler 前从上下文提取属性
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	// 从上下文中提取预存的 slog 属性
	if attrs, ok := ctx.Value("slog_attrs").([]slog.Attr); ok {
		for _, attr := range attrs {
			r.AddAttrs(attr) // 将上下文属性添加到日志记录
		}
	}

	// 调用底层 Handler 的 Handle 方法
	return h.Handler.Handle(ctx, r)
}

// WithLogAttr 辅助函数，用于向上下文添加 slog 属性
func WithLogAttr(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	// 获取现有的属性列表，或创建新的
	if existing, ok := parent.Value("slog_attrs").([]slog.Attr); ok {
		existing = append(existing, attr)
		return context.WithValue(parent, "slog_attrs", existing)
	}

	// 首次添加属性
	return context.WithValue(parent, "slog_attrs", []slog.Attr{attr})
}

func main() {
	// 创建 ContextHandler 包装 JSONHandler
	handler := &ContextHandler{
		Handler: slog.NewJSONHandler(os.Stdout, nil),
	}
	logger := slog.New(handler)

	// 创建包含请求信息的上下文
	ctx := WithLogAttr(context.Background(),
		slog.String("request_id", "req-12345"),
	)
	ctx = WithLogAttr(ctx, slog.String("user_id", "jianghushinian"))

	// 使用带上下文的日志方法
	logger.InfoContext(ctx, "User login successful",
		slog.String("action", "login"),
	)
}
