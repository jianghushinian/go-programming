//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
)

// func main() {
// 	onceBody := sync.OnceFunc(func() {
// 		fmt.Println("Only once")
// 	})
// 	done := make(chan bool)
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			onceBody()
// 			done <- true
// 		}()
// 	}
// 	for i := 0; i < 10; i++ {
// 		<-done
// 	}
// }

func main() {
	onceBody := sync.OnceFunc(func() {
		panic("Only once")
	})

	for i := 0; i < 3; i++ {
		func() {
			defer func() {
				r := recover()
				fmt.Println("recover", r)
			}()
			onceBody()
		}()
	}
}
