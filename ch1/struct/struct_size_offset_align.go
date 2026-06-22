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
	t1 := T1{}
	fmt.Println("# T1")
	fmt.Printf("T1.a: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t1.a), unsafe.Offsetof(t1.a), unsafe.Alignof(t1.a))
	fmt.Printf("T1.b: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t1.b), unsafe.Offsetof(t1.b), unsafe.Alignof(t1.b))
	fmt.Printf("T1.c: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t1.c), unsafe.Offsetof(t1.c), unsafe.Alignof(t1.c))

	t2 := T2{}
	fmt.Println("# T2")
	fmt.Printf("T2.b: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t2.b), unsafe.Offsetof(t2.b), unsafe.Alignof(t2.b))
	fmt.Printf("T2.a: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t2.a), unsafe.Offsetof(t2.a), unsafe.Alignof(t2.a))
	fmt.Printf("T2.c: size=%d, offset=%v, align=%d\n", unsafe.Sizeof(t2.c), unsafe.Offsetof(t2.c), unsafe.Alignof(t2.c))
}
