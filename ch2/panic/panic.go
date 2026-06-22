//go:build ignore
// +build ignore

package main

import "fmt"

func f() {
	defer fmt.Println("defer 1")
	fmt.Println(1)
	panic("woah")
	defer fmt.Println("defer 2")
	fmt.Println(2)
}

func main() {
	f()
}
