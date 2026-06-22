//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	{
		type Empty struct{}

		var s1 struct{}
		s2 := Empty{}
		s3 := struct{}{}

		fmt.Printf("s1 addr: %p, size: %d\n", &s1, unsafe.Sizeof(s1))
		fmt.Printf("s2 addr: %p, size: %d\n", &s2, unsafe.Sizeof(s2))
		fmt.Printf("s3 addr: %p, size: %d\n", &s3, unsafe.Sizeof(s3))
		fmt.Printf("s1 == s2 == s3: %t\n", s1 == s2 && s2 == s3)
	}

	{
		var (
			a struct{}
			b struct{}
			c struct{}
			d struct{}
		)

		println("&a:", &a)
		println("&b:", &b)
		println("&c:", &c)
		println("&d:", &d)

		println("&a == &b:", &a == &b)
		x := &a
		y := &b
		println("x == y:", x == y)

		fmt.Printf("&c(%p) == &d(%p): %t\n", &c, &d, &c == &d)
	}
}
