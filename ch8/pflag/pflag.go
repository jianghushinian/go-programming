//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

func main() {
	// 1. <Type>：返回指针
	name := pflag.String("name", "guest", "user name")

	// 2. <Type>Var：将标志值绑定到已有变量
	var age int
	pflag.IntVar(&age, "age", 18, "user age")

	// 3. <Type>P / <Type>VarP：支持短标志
	var city string
	pflag.StringVarP(&city, "city", "c", "Shanghai", "living city")

	// 解析命令行参数
	pflag.Parse()

	fmt.Println("Name:", *name)
	fmt.Println("Age:", age)
	fmt.Println("City:", city)
}
