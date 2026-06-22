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

func main() {
	t1 := T1{}
	fmt.Printf("T1: %+v\n", t1)
	b := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&t1)) + unsafe.Offsetof(t1.b)))
	*b = "江湖十年"
	fmt.Printf("T1: %+v\n", t1)
}
