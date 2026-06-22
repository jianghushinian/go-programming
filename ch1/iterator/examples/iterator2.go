//go:build ignore
// +build ignore

package main

func iterators(slice []int) func(yield func(i, v int) bool) {
	return func(yield func(i int, v int) bool) {
		for i, v := range slice {
			if !yield(i, v) {
				return
			}
		}
	}
}

func iteratorm(m map[string]int) func(yield func(k string, v int) bool) {
	return func(yield func(k string, v int) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}
