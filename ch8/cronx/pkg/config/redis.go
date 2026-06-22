package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig defines config for mysql database.
type RedisConfig struct {
	Addr         string
	Username     string
	Password     string
	Database     int
	MaxRetries   int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
	PoolSize     int
}

// NewRedis create a new gorm db instance with the given config.
func NewRedis(cfg *RedisConfig) (*redis.Client, error) {
	config := &redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.Database,
		MaxRetries:   cfg.MaxRetries,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  cfg.PoolTimeout,
		PoolSize:     cfg.PoolSize,
	}

	rdb := redis.NewClient(config)

	// check redis if is ok.
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
