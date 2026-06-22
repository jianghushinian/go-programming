//go:build ignore
// +build ignore

package singleton

type Config struct {
	Name string
}

var instance = &Config{Name: "江湖十年"} // 程序启动时直接初始化

func GetInstance() *Config {
	return instance
}
