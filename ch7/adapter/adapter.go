package adapter

import (
	"log"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 统一日志接口
type Logger interface {
	Debug(msg string, fields map[string]any)
	Info(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
}

// ZapAdapter zap 日志库适配器
type ZapAdapter struct {
	logger *zap.Logger // 持有被适配对象
}

func NewZapAdapter(logger *zap.Logger) *ZapAdapter {
	return &ZapAdapter{logger: logger}
}

func (z *ZapAdapter) Debug(msg string, fields map[string]any) {
	z.logger.Debug(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Info(msg string, fields map[string]any) {
	z.logger.Info(msg, convertFields(fields)...)
}

func (z *ZapAdapter) Error(msg string, fields map[string]any) {
	z.logger.Error(msg, convertFields(fields)...)
}

func convertFields(fields map[string]any) []zapcore.Field {
	if len(fields) == 0 {
		return nil
	}
	result := make([]zapcore.Field, 0, len(fields))
	for k, v := range fields {
		result = append(result, zap.Any(k, v))
	}
	return result
}

// LogrusAdapter Logrus 日志库适配器
type LogrusAdapter struct {
	logger *logrus.Logger
}

func NewLogrusAdapter(logger *logrus.Logger) *LogrusAdapter {
	return &LogrusAdapter{logger: logger}
}

func (l *LogrusAdapter) Debug(msg string, fields map[string]any) {
	l.logger.WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l *LogrusAdapter) Info(msg string, fields map[string]any) {
	l.logger.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *LogrusAdapter) Error(msg string, fields map[string]any) {
	l.logger.WithFields(logrus.Fields(fields)).Error(msg)
}

// StdLibAdapter 标准库日志适配器
type StdLibAdapter struct{}

func NewStdLibAdapter() *StdLibAdapter {
	return &StdLibAdapter{}
}

func (s *StdLibAdapter) Debug(msg string, fields map[string]any) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

func (s *StdLibAdapter) Info(msg string, fields map[string]any) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (s *StdLibAdapter) Error(msg string, fields map[string]any) {
	log.Printf("[ERROR] %s %v", msg, fields)
}
