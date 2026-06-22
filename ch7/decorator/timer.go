package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func Timeit(fn func()) func() {
	return func() {
		start := time.Now()
		fn()
		fmt.Printf("%s took %.4fs\n", getFunctionName(fn), time.Since(start).Seconds())
	}
}

func foo() {
	time.Sleep(1 * time.Second)
	fmt.Println("done")
}

func getFunctionName(fn any) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	// 提取最后的函数名部分，示例格式：main.foo
	parts := strings.Split(fullName, ".")
	return parts[len(parts)-1]
}

func main() {
	decorated := Timeit(foo)
	decorated()
}
