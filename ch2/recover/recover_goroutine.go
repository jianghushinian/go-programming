//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"
)

func f() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover:", r)
		}
	}()

	go func() {
		panic("woah")
	}()
	time.Sleep(1 * time.Second)
}

func main() {
	f()
}
