//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"unsafe"
)

type T3 struct {
	a int8
	b string
	c struct{}
}

type T4 struct {
	a int8
	c struct{}
	b string
}

type T5 struct {
	c struct{}
	a int8
	b string
}

func main() {
	fmt.Printf("T3 size: %d\n", unsafe.Sizeof(T3{}))
	fmt.Printf("T4 size: %d\n", unsafe.Sizeof(T4{}))
	fmt.Printf("T5 size: %d\n", unsafe.Sizeof(T5{}))

	t3 := T3{}
	fmt.Println("# T3")
	fmt.Printf("T3.a: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t3.a), unsafe.Offsetof(t3.a), unsafe.Alignof(t3.a))
	fmt.Printf("T3.b: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t3.b), unsafe.Offsetof(t3.b), unsafe.Alignof(t3.b))
	fmt.Printf("T3.c: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t3.c), unsafe.Offsetof(t3.c), unsafe.Alignof(t3.c))
}
