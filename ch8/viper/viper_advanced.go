//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	viper.Set("database.host", "127.0.0.1")
	fmt.Println(viper.GetString("database.host")) // 127.0.0.1

	pflag.String("app.mode", "dev", "运行模式")
	pflag.Int("app.port", 8080, "服务端口")
	pflag.String("app.name", "myapp", "应用名")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)
	viper.BindPFlag("app.name", pflag.Lookup("app.name"))

	fmt.Println("Mode:", viper.GetString("app.mode"))
	fmt.Println("Port:", viper.GetInt("app.port"))
	fmt.Println("Name:", viper.GetString("app.name"))

	// export DB_USER=admin
	// export DB_PASS=123456
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.pass", "DB_PASS")
	fmt.Println(viper.GetString("database.user")) // admin
	fmt.Println(viper.GetString("database.pass")) // 123456

	viper.AutomaticEnv()
	fmt.Println(viper.Get("PATH")) // 自动读取系统环境变量 PATH

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	fmt.Println(viper.Get("server.ip")) // 输出 127.0.0.1
	fmt.Println(viper.Get("server-ip")) // 同样输出 127.0.0.1

	viper.AllowEmptyEnv(true)

	viper.SetEnvPrefix("MYAPP")
	viper.BindEnv("app") // 对应环境变量为 MYAPP_APP
	fmt.Println(viper.Get("APP"))

	var yamlExample = []byte(`
server:
  address: 127.0.0.1
  port: 8080
  timeout: 10
`)
	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(yamlExample))
	fmt.Println(viper.Get("server.timeout"))

	viper.SetDefault("app.name", "myapp")
	viper.SetDefault("app.port", 8080)

	fmt.Println(viper.GetString("app.name")) // 输出: myapp
	fmt.Println(viper.GetInt("app.port"))    // 输出: 8080

	fmt.Println(viper.Get("app.name"))
	fmt.Println(viper.GetString("app.name"))

	fmt.Println(viper.AllKeys())     // 返回所有键的切片
	fmt.Println(viper.AllSettings()) // 返回 map[string]any

	if viper.IsSet("database.user") {
		fmt.Println("database user already exists")
	}

	appconf := viper.Sub("app")      // 提取 app 子树
	fmt.Println(appconf.Get("name")) // 返回 myapp

}
