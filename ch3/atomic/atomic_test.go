package main

import (
	"sync/atomic"
	"testing"
)

func TestAtomicAdd(t *testing.T) {
	i32 := int32(4)
	atomic.AddInt32(&i32, 2)
	if i32 != 6 {
		t.Fatalf("expected 6, got %d", i32)
	}
}

func TestAtomicLoad(t *testing.T) {
	i32 := int32(42)
	if v := atomic.LoadInt32(&i32); v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestAtomicStore(t *testing.T) {
	i32 := int32(10)
	atomic.StoreInt32(&i32, 42)
	if i32 != 42 {
		t.Fatalf("expected 42, got %d", i32)
	}
}

func TestAtomicSwap(t *testing.T) {
	i32 := int32(10)
	if v := atomic.SwapInt32(&i32, 42); v != 10 {
		t.Fatalf("expected old value 10, got %d", v)
	}
	if i32 != 42 {
		t.Fatalf("expected new value 42, got %d", i32)
	}
}

func TestAtomicCompareAndSwap(t *testing.T) {
	i32 := int32(10)
	if !atomic.CompareAndSwapInt32(&i32, 10, 42) {
		t.Fatalf("expected CAS to succeed")
	}
	if i32 != 42 {
		t.Fatalf("expected new value 42, got %d", i32)
	}
	if atomic.CompareAndSwapInt32(&i32, 10, 100) {
		t.Fatalf("expected CAS to fail")
	}
	if i32 != 42 {
		t.Fatalf("expected value to remain 42, got %d", i32)
	}
}

func TestAtomicValue(t *testing.T) {
	var vs atomic.Value
	vs.Store("hello")
	vs.Swap("world")
	vs.CompareAndSwap("world", "atomic")
	if vs.Load().(string) != "atomic" {
		t.Fatalf("expected 'atomic', got %s", vs.Load().(string))
	}
}

func TestAtomicInt32(t *testing.T) {
	var i32 atomic.Int32
	i32.Add(4)
	if v := i32.Load(); v != 4 {
		t.Fatalf("expected 4, got %d", v)
	}
}
