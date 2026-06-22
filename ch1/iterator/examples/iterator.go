//go:build ignore
// +build ignore

package main

import "fmt"

/*
func iterator(yield func() bool) {
	for i := 0; i < 5; i++ {
		if !yield() {
			return
		}
	}
}

func main() {
	i := 0
	for range iterator {
		fmt.Printf("i=%d\n", i)
		i++
	}
}
*/

func iterator(n int) func(yield func() bool) {
	return func(yield func() bool) {
		for i := 0; i < n; i++ {
			if !yield() {
				return
			}
		}
	}
}

func main() {
	i := 0
	for range iterator(3) {
		fmt.Printf("i=%d\n", i)
		i++
	}
}
