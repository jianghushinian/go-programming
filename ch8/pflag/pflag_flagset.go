//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	// 创建两个子命令：add 和 delete
	addCmd := pflag.NewFlagSet("add", pflag.ExitOnError)
	delCmd := pflag.NewFlagSet("delete", pflag.ExitOnError)

	// 为 add 子命令定义标志
	addName := addCmd.String("name", "", "name to add")

	// 为 delete 子命令定义标志
	delID := delCmd.Int("id", 0, "id to delete")

	switch os.Args[1] {
	case "add":
		_ = addCmd.Parse(os.Args[2:]) // 解析 add 子命令的参数
		fmt.Printf("Add command called: name=%s\n", *addName)
	case "delete":
		_ = delCmd.Parse(os.Args[2:]) // 解析 delete 子命令的参数
		fmt.Printf("Delete command called: id=%d\n", *delID)
	default:
		fmt.Println("expected 'add' or 'delete' subcommands")
		os.Exit(1)
	}
}
