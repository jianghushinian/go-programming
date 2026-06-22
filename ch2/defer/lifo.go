//go:build ignore
// +build ignore

package main

import "fmt"

func f1() {
	defer fmt.Println("deferred in f 1")
	defer fmt.Println("deferred in f 2")
	defer fmt.Println("deferred in f 3")
	fmt.Println("calling f")
}

func f2() {
	fmt.Println("1")

	defer func() {
		fmt.Println("2")
		defer fmt.Println("3")
		fmt.Println("4")
	}()

	fmt.Println("5")

	defer fmt.Println("6")

	fmt.Println("7")
}

func main() {
	f1()
	f2()
}
