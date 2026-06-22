//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

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

	// 获取配置项
	fmt.Println("App Name:", viper.GetString("app.name"))
	fmt.Println("App Mode:", viper.GetString("app.mode"))

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config %s changed\n", e.Name)
	})

	viper.Set("app.mode", "prod")
	err := viper.WriteConfig() // 写回到当前配置文件
	if err != nil {
		panic(err)
	}
}
