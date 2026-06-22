package main

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", k)
			m.Store(key, k)
			m.Load(key)
		}(i % 100000) // 使用 100000 个不同的 key
	}
	wg.Wait()
}

func BenchmarkRWMutexMap(b *testing.B) {
	var m = NewRWMutexMap()
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", k)
			m.Set(key, k)
			m.Get(key)
		}(i % 100000)
	}
	wg.Wait()
}

func BenchmarkShardingMap(b *testing.B) {
	var m = NewShardingMap(10)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", k)
			m.Set(key, k)
			m.Get(key)
		}(i % 100000)
	}
	wg.Wait()
}
