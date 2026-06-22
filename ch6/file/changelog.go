package main

import (
	"io"
	"os"
)

var (
	version       = "dev"          // 版本号
	changeLogPath = "CHANGELOG.md" // 文件路径
)

type ChangeLogSpec struct {
	Version   string
	ChangeLog string
}

// GetChangeLog 读取 ChangeLog 文件并返回完整信息
func GetChangeLog() (ChangeLogSpec, error) {
	data, err := os.ReadFile(changeLogPath)
	if err != nil {
		return ChangeLogSpec{}, err
	}
	return ChangeLogSpec{
		Version:   version,
		ChangeLog: string(data),
	}, nil
}

func GetChangeLogReader(r io.Reader) (ChangeLogSpec, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return ChangeLogSpec{}, err
	}
	return ChangeLogSpec{
		Version:   version,
		ChangeLog: string(data),
	}, nil
}
