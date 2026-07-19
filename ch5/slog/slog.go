//go:build ignore
// +build ignore

package main

import (
	"context"
	"log"
	"log/slog"
	"os"
)

func main() {
	slog.Debug("debug message")
	slog.Info("info message", "name", "江湖十年", "age", 20)
	slog.WarnContext(context.Background(), "warn message")
	slog.Error("error message")

	slog.SetLogLoggerLevel(slog.LevelDebug)

	jsonLogger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,            // 记录日志产生位置
			Level:     slog.LevelDebug, // 设置日志级别
		}),
	)
	jsonLogger.Info("info message", "name", "江湖十年", "age", 20)

	textLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	textLogger.Info("info message", "name", "江湖十年", "age", 20)

	slog.SetDefault(jsonLogger)
	slog.Info("info message", "name", "江湖十年", "age", 20)

	log.Println("standard log")

	jsonLogger.Info(
		"info message",
		slog.Group("user", "name", "江湖十年", slog.Int("age", 20)),
	)
	textLogger.Info(
		"info message",
		slog.Group("user", "name", "江湖十年", slog.Int("age", 20)),
	)

	childLogger := textLogger.With(
		slog.String("request_id", "req-12345"),
		slog.String("service", "auth"),
	)
	childLogger.Info("user authenticated", slog.String("user_id", "jianghushinian"))

	user := User{ID: "jianghushinian", Name: "江湖十年", Password: "secret"}
	textLogger.Info("user login", "user", user)

	const (
		LevelTrace = slog.Level(-2)
		LevelFatal = slog.Level(12)
	)

	// 在 HandlerOptions 中自定义级别名称
	opts := &slog.HandlerOptions{
		Level: LevelTrace,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey { // 日志级别
				level := a.Value.Any().(slog.Level)
				switch level {
				case LevelTrace: // 自定义级别
					a.Value = slog.StringValue("TRACE")
				case LevelFatal: // 自定义级别
					a.Value = slog.StringValue("FATAL")
				}
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	logger.LogAttrs(context.Background(), LevelTrace, "trace message", slog.String("name", "江湖十年"))
	logger.LogAttrs(context.Background(), LevelFatal, "fatal message", slog.String("name", "江湖十年"))

	// 输出至多个 Handler
	file, _ := os.OpenFile("app.jsonl", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	multiHandler := slog.NewMultiHandler(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	multiLogger := slog.New(multiHandler)
	multiLogger.Debug("Text Handler log", "request_id", "req-12345")
	multiLogger.Info("Multi Handler log", "user_id", "jianghushinian")
}

type User struct {
	ID       string
	Name     string
	Password string
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", u.ID),
		slog.String("name", u.Name),
		// Password 字段被隐藏
	)
}
