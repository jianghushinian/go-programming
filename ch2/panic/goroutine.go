//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"
)

func g() {
	fmt.Println("calling g")
	go f(0)
	fmt.Println("called g")
}

func f(i int) {
	fmt.Println("panicking!")
	panic(fmt.Sprintf("i=%v", i))
	fmt.Println("printing in f", i)
}

func main() {
	g()
	time.Sleep(10 * time.Second)
}
