//go:build ignore
// +build ignore

package main

import "fmt"

func f1() {
	defer func() {
		recover()
	}()

	defer fmt.Println("defer 1")
	fmt.Println(1)
	panic("woah")
	defer fmt.Println("defer 2")
	fmt.Println(2)
}

func f2() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover:", r)
		}
	}()
	panic("woah")
}

func main() {
	// f1()
	f2()
}
