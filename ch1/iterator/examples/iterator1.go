//go:build ignore
// +build ignore

package main

import "fmt"

func iterator1(n int) func(yield func(v int) bool) {
	return func(yield func(v int) bool) {
		for i := 0; i < n; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	i := 0
	for v := range iterator1(10) {
		if i >= 5 {
			break
		}
		fmt.Printf("%d => %d\n", i, v)
		i++
	}
}
