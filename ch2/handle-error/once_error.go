//go:build ignore
// +build ignore

package main

import (
	"errors"
	"fmt"
	"log"
)

// func MustFoo() string {
// 	data, err := Foo()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return data
// }

func Foo() error {
	return errors.New("foo error")
}

func Bar1() error {
	return Foo()
}

func Bar2() error {
	err := Foo()
	if err != nil {
		// 提供服务降级处理，如记录日志
		log.Println("Bar2 error:", err)
		return nil
	}
	// ...
	return nil
}

func main() {
	err := Bar1()
	if err != nil {
		// do something
		fmt.Println("bar error:", err)
	}

	err = Bar2()
	if err != nil {
		// do something
		fmt.Println("bar error:", err)
	}
}
