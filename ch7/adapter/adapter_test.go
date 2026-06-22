package adapter

import (
	"testing"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func TestAdapter(t *testing.T) {
	zapLogger, _ := zap.NewProduction()
	logrusLogger := logrus.New()

	var logger Logger

	// 动态切换日志实现
	logger = NewZapAdapter(zapLogger)
	logger.Debug("From zap", map[string]any{"name": "江湖十年"})

	logger = NewLogrusAdapter(logrusLogger)
	logger.Debug("From logrus", map[string]any{"name": "江湖十年"})

	logger = NewStdLibAdapter()
	logger.Debug("From log", map[string]any{"name": "江湖十年"})
}
