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
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = &Config{Name: "江湖十年"}
	}
	return instance
}
