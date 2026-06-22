//go:build ignore
// +build ignore

package main

import "fmt"

func f1() (r int) {
	r = 2
	defer func() {
		fmt.Println("r:", r)
		r *= 3
	}()
	return r
}

func f2() (r int) {
	defer func(r int) {
		fmt.Println("r:", r)
		r *= 3
	}(r)
	return 2
}

func main() {
	fmt.Println(f1())
	fmt.Println(f2())
}
