package main

import (
	"embed"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChangeLog(t *testing.T) {
	// 创建临时文件
	f, err := os.CreateTemp("", "TEST_CHANGELOG")
	assert.NoError(t, err)
	defer func() {
		_ = f.Close()
		_ = os.RemoveAll(f.Name()) // 测试完成后清理
	}()

	// 准备测试数据
	testData := `# Changelog
## v1.0.0
- Initial release`
	_, err = f.WriteString(testData)
	assert.NoError(t, err)

	// 临时替换文件路径和全局变量
	version = "v1.0.0"
	changeLogPath = f.Name()

	// 执行测试
	actual, err := GetChangeLog()
	assert.NoError(t, err)
	expected := ChangeLogSpec{Version: "v1.0.0", ChangeLog: testData}
	assert.Equal(t, expected, actual)
}

//go:embed testdata/*.md
var testDataFS embed.FS

func TestGetChangeLogReader(t *testing.T) {
	// 从嵌入的文件系统中读取内容
	testData, err := testDataFS.ReadFile("testdata/changelog.md")
	assert.NoError(t, err)

	// 设置测试所需的全局变量值
	originalVersion := version
	version = "v1.0.0"
	// 测试完成后恢复全局变量
	defer func() { version = originalVersion }()
	// 创建字符串读取器作为输入
	reader := strings.NewReader(string(testData))

	// 执行测试
	actual, err := GetChangeLogReader(reader)
	assert.NoError(t, err)
	expected := ChangeLogSpec{Version: "v1.0.0", ChangeLog: string(testData)}
	assert.Equal(t, expected, actual)
}
