package singleton

import "sync"

type Config struct {
	Name string
}

var (
	instance *Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{Name: "江湖十年"}
	})
	return instance
}
