//go:build ignore
// +build ignore

package singleton

import "sync"

type Config struct {
	Name string
}

var (
	instance *Config
	mu       sync.Mutex
)

func GetInstance() *Config {
	if instance == nil { // 第一次检查（无锁）
		mu.Lock() // 加锁
		defer mu.Unlock()
		if instance == nil { // 第二次检查（加锁后）
			instance = &Config{Name: "江湖十年"}
		}
	}
	return instance
}
