//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"reflect"
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
	for _, T := range []any{T1{}, T2{}} {
		typ := reflect.TypeOf(T)
		fmt.Printf("%s size: %d\n", typ.Name(), typ.Size())

		n := typ.NumField()
		for i := 0; i < n; i++ {
			field := typ.Field(i)
			fmt.Printf("%s.%s: size=%d, offset=%v, align=%d\n",
				typ.Name(),
				field.Name,
				field.Type.Size(),
				field.Offset,
				field.Type.Align(),
			)
		}
	}
}
