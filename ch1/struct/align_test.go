package main

import "testing"

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

func Benchmark_T1_Align(b *testing.B) {
	b.ReportAllocs() // 开启内存统计
	b.ResetTimer()   // 重置计时器
	for i := 0; i < b.N; i++ {
		_ = make([]T1, b.N)
	}
}

func Benchmark_T2_Align(b *testing.B) {
	b.ReportAllocs() // 开启内存统计
	b.ResetTimer()   // 重置计时器
	for i := 0; i < b.N; i++ {
		_ = make([]T2, b.N)
	}
}
