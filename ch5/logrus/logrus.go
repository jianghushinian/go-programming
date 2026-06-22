//go:build ignore
// +build ignore

package main

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel) // 日志级别设为 Debug
	log.Debug("Init")
	log.Info("Server started")
	log.Warn("Low disk space")
	log.Error("Failed to connect to database")

	log.SetFormatter(&log.JSONFormatter{})
	log.WithField("module", "payment").Info("payment processed")

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true, // 增加时间戳
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true, // 关闭颜色
	})
	log.WithField("module", "payment").Info("payment processed")

	// log.SetFormatter(&log.TextFormatter{})
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(log.Fields{
		"user_id": 42,
		"action":  "login",
		"ip":      "192.168.0.1",
	}).Info("user login success")

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Warn("Failed to log to file, using default stderr")
	} else {
		// log.SetOutput(file)
		mw := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mw)
	}
	log.Info("Log to multi writer")

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetReportCaller(true) // 记录文件名和函数名
	reqLogger := logger.WithField("request_id", "xxx")
	reqLogger.Info("something happened on that request") // 会记录 request_id
}
