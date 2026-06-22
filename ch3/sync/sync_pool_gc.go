//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	count := 0
	pool := &sync.Pool{
		New: func() any {
			count++
			return fmt.Sprintf("obj(%d)", count)
		},
	}

	obj := pool.Get()
	fmt.Printf("第 1 次获取对象: %v\n", obj)
	pool.Put(obj)

	obj = pool.Get()
	fmt.Printf("第 2 次获取对象: %v\n", obj)
	pool.Put(obj)

	runtime.GC() // 连续两次 GC
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	obj = pool.Get()
	fmt.Printf("第 3 次获取对象: %v\n", obj)
}
