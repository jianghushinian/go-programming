//go:build ignore
// +build ignore

package main

import "fmt"

func deferNil() {
	var f func()
	defer f()
	fmt.Println("calling deferNil")
}

func main() {
	deferNil()
}
