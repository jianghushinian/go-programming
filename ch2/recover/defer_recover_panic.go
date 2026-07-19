//go:build ignore
// +build ignore

package main

import "fmt"

func f() {
	defer func() {
		panic("woah 1")
	}()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover:", r)
		}
	}()

	panic("woah 2")
}

func main() {
	f()
}
