//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	s := []string{"a", "b", "c"}
	for i, v := range slices.All(s) {
		fmt.Printf("%d => %s\n", i, v)
	}

	m := map[string]int{"a": 0, "b": 1, "c": 2}
	for k, v := range maps.All(m) {
		fmt.Printf("%s: %d\n", k, v)
	}
}
