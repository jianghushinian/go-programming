//go:build ignore
// +build ignore

package main

import "fmt"

func main() {
	m := make(map[string]int)

	go func() {
		for {
			m["k"] = 1
			fmt.Println("set k:", 1)
		}
	}()

	go func() {
		for {
			v, _ := m["k"]
			fmt.Println("read k:", v)
		}
	}()

	select {}
}
