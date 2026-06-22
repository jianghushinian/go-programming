package main

import (
	"flag"
	"fmt"
)

type nameValue string

func (n *nameValue) String() string {
	return string(*n)
}

func (n *nameValue) Set(s string) error {
	*n = nameValue(s)
	return nil
}

func main() {
	// 定义命令行参数选项
	var nFlag = flag.Int("n", 1234, "help message for flag n")
	var flagvar int
	flag.IntVar(&flagvar, "flagname", 5678, "help message for flagname")
	var flagVal nameValue
	flag.Var(&flagVal, "name", "help message for name")

	// 解析命令行参数
	flag.Parse()

	fmt.Println("n has value:", *nFlag)
	fmt.Println("flagvar has value:", flagvar)
	fmt.Println("name has value:", flagVal)
}
