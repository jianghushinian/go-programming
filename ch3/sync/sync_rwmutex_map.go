package main

import "sync"

type RWMutexMap struct {
	rw sync.RWMutex
	m  map[string]int
}

func NewRWMutexMap() *RWMutexMap {
	return &RWMutexMap{
		m: make(map[string]int),
	}
}

func (m *RWMutexMap) Set(key string, v int) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.m[key] = v
}

func (m *RWMutexMap) Get(key string) (int, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	v, ok := m.m[key]
	return v, ok
}

func (m *RWMutexMap) Del(key string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	delete(m.m, key)
}
