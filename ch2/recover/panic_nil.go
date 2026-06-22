//go:build ignore
// +build ignore

package main

import "fmt"

func f() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic(nil)
}

func main() {
	f()
}
