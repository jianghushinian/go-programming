//go:build ignore
// +build ignore

package main

import "fmt"

/*
func iterator(slice []int) func(yield func(i, v int) bool) {
	return func(yield func(i int, v int) bool) {
		for i, v := range slice {
			if !yield(i, v) {
				return
			}
		}
	}
}

func main() {
	s := []int{1, 2, 3, 4, 5}
	for i, v := range iterator(s) {
		if i == 3 {
			break
		}
		fmt.Printf("%d => %d\n", i, v)
	}
}
*/

func iterator(slice []int) func(yield func(i, v int) bool) {
	return func(yield func(i int, v int) bool) {
		for i, v := range slice {
			if !yield(i, v) {
				return
			}
		}
	}
}

func yield(i, v int) bool {
	if i == 3 {
		return false
	}
	fmt.Printf("%d => %d\n", i, v)
	return true
}

func main() {
	s := []int{1, 2, 3, 4, 5}
	iterator(s)(yield)
}
