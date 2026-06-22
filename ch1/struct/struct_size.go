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
	fmt.Printf("T1 size: %d\n", unsafe.Sizeof(T1{}))
	fmt.Printf("T2 size: %d\n", unsafe.Sizeof(T2{}))
}
