//go:build ignore
// +build ignore

package main

import "sync/atomic"

// type S1 struct {
// 	a int32
// 	b int64
// }

type S1 struct {
	a   int32
	pad uint32 // ensure 8-byte alignment of val on 386
	b   int64
}

func main() {
	s1 := S1{}
	atomic.AddInt64(&s1.b, 1)
}
