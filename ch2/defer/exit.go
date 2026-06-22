//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
)

func f() {
	defer fmt.Println("deferred in f")
	fmt.Println("calling f")
	os.Exit(0)
}

func main() {
	f()
}
