//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

func main() {
	var ip = pflag.IntP("flagname", "f", 1234, "help message")
	pflag.Lookup("flagname").NoOptDefVal = "4321"

	var logLevel = pflag.String("log-level", "INFO", "help message")
	pflag.CommandLine.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
	})

	config := pflag.String("config", "", "config file")
	conf := pflag.StringP("conf", "c", "", "deprecated, use --config instead")
	// pflag.CommandLine.MarkDeprecated("conf", "please use --config")
	pflag.CommandLine.MarkShorthandDeprecated("c", "please use --config")
	pflag.CommandLine.MarkHidden("conf")

	// 解析命令行参数
	pflag.Parse()

	fmt.Println("IP:", *ip)
	fmt.Println("Log Level:", *logLevel)
	fmt.Println("Config:", *config)
	fmt.Println("Conf:", *conf)
}
