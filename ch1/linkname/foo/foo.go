package foo

import (
	_ "unsafe"

	// 被拉取的包需要显式导入
	_ "linkname/bar"
)

//go:linkname Add linkname/bar.add
func Add(a, b int) int

func Div(a, b int) int

//go:linkname Hello linkname/bar.hello
func Hello(name string) string