package main

import (
	"hash/maphash"
	"sync"
)

var seed = maphash.MakeSeed()

func hashKey(key string) uint64 {
	return maphash.String(seed, key)
}

type ShardingMap struct {
	locks  []sync.RWMutex
	shards []map[string]int
}

func NewShardingMap(size int) *ShardingMap {
	sm := &ShardingMap{
		locks:  make([]sync.RWMutex, size),
		shards: make([]map[string]int, size),
	}
	for i := 0; i < size; i++ {
		sm.shards[i] = make(map[string]int)
	}
	return sm
}

func (m *ShardingMap) getShardIdx(key string) uint64 {
	hash := hashKey(key)
	return hash % uint64(len(m.shards))
}

func (m *ShardingMap) Set(key string, value int) {
	idx := m.getShardIdx(key)
	m.locks[idx].Lock()
	defer m.locks[idx].Unlock()
	m.shards[idx][key] = value
}

func (m *ShardingMap) Get(key string) (int, bool) {
	idx := m.getShardIdx(key)
	m.locks[idx].RLock()
	defer m.locks[idx].RUnlock()
	value, ok := m.shards[idx][key]
	return value, ok
}

func (m *ShardingMap) Del(key string) {
	idx := m.getShardIdx(key)
	m.locks[idx].Lock()
	defer m.locks[idx].Unlock()
	delete(m.shards[idx], key)
}
