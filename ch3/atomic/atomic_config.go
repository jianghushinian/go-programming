package main

import (
	"fmt"
	"sync/atomic"
)

type Config struct {
	Name string
	Age  int
	Addr string
}

type AtomicConfig struct {
	config atomic.Value // 存储 *Config
}

func NewAtomicConfig(cfg *Config) *AtomicConfig {
	c := &AtomicConfig{}
	c.config.Store(cfg)
	return c
}

func (c *AtomicConfig) Set(newCfg *Config) {
	c.config.Store(newCfg) // 内部保证原子性
}

func (c *AtomicConfig) Get() *Config {
	return c.config.Load().(*Config) // 类型断言
}

func main() {
	cfg := &Config{Name: "江湖十年", Age: 20, Addr: "杭州"}
	atomicConfig := NewAtomicConfig(cfg)
	fmt.Printf("Config inited: %+v\n", atomicConfig.Get())

	newCfg := &Config{Name: "jianghushinian", Age: 30, Addr: "ShangHai"}
	atomicConfig.Set(newCfg)
	fmt.Printf("Config updated: %+v\n", atomicConfig.Get())
}
