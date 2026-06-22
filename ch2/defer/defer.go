//go:build ignore
// +build ignore

package main

import "fmt"

func f() {
	defer fmt.Println("deferred in f")
	fmt.Println("calling f")
}

func main() {
	f()
}
