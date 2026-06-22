//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"unsafe"
)

type T1 struct {
	a int8
	b string
	c bool
}

type T2 struct {
	b string
	a int8
	c bool
}

func main() {
	fmt.Printf("int8 align: %d\n", unsafe.Alignof(int8(1)))
	fmt.Printf("bool align: %d\n", unsafe.Alignof(true))
	fmt.Printf("string align: %d\n", unsafe.Alignof("Hello World"))
	fmt.Printf("T1 align: %d\n", unsafe.Alignof(T1{}))
	fmt.Printf("T2 align: %d\n", unsafe.Alignof(T1{}))
	fmt.Printf("empty struct align: %d\n", unsafe.Alignof(struct{}{}))
	fmt.Printf("int align: %d\n", unsafe.Alignof(int(3)))
	fmt.Printf("int array align: %d\n", unsafe.Alignof([3]int{1, 2, 3}))
}
