//go:build ignore
// +build ignore

package main

import "fmt"

func f() {
	m := map[int]struct{}{}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("goroutine 1", err)
			}
		}()
		for {
			m[1] = struct{}{}
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("goroutine 2", err)
			}
		}()
		for {
			m[1] = struct{}{}
		}
	}()

	select {}
}

func main() {
	f()
}
