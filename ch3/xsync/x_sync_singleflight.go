//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type User struct {
	Id   int64
	Name string
}

func GetUserFromDB(username string) *User {
	fmt.Printf("Querying DB for key: %s\n", username)
	time.Sleep(1 * time.Second) // 模拟耗时
	id, _ := strconv.Atoi(username[len(username)-3:])
	return &User{
		Id:   int64(id),
		Name: username,
	}
}

var (
	cache    sync.Map // 缓存
	reqGroup singleflight.Group
)

func GetUser(key string) *User {
	// 先尝试从缓存获取
	val, ok := cache.Load(key)
	if ok {
		return val.(*User)
	}

	fmt.Printf("User %s not in cache\n", key)

	// 缓存未命中，使用 singleFlight.Group 合并请求
	result, _, _ := reqGroup.Do(key, func() (any, error) {
		// 从数据库读取数据
		val := GetUserFromDB(key)
		// 存入缓存
		cache.Store(key, val)
		return val, nil
	})
	return result.(*User)
}

func main() {
	var wg sync.WaitGroup
	keys := []string{"user_123", "user_123", "user_123"}

	for _, key := range keys {
		wg.Go(func() {
			fmt.Printf("Get user %s: %+v\n", key, GetUser(key))
		})
	}
	wg.Wait()
}
