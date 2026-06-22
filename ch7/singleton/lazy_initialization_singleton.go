//go:build ignore
// +build ignore

package singleton

type Config struct {
	Name string
}

var instance *Config

func GetInstance() *Config {
	if instance == nil {
		instance = &Config{Name: "江湖十年"} // 第一次访问时才创建
	}
	return instance
}
