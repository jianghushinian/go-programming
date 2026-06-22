//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Port int
	}
	Database struct {
		User string
		Pass string
		Host string
	}
}

func main() {
	// 设置配置文件名（不带扩展名）
	viper.SetConfigName("config")
	// 指定文件类型
	viper.SetConfigType("yaml")
	// 设置配置文件路径（可多次调用）
	viper.AddConfigPath(".")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config error: %v", err))
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	fmt.Println(cfg.App.Name)
	fmt.Println(cfg.Database.Host)
}
