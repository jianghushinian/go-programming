package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

func main() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     "localhost:6379",
		Password: "password",
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	mutex := rs.NewMutex("test-redsync",
		redsync.WithExpiry(30*time.Second),
		redsync.WithTries(3),
		redsync.WithRetryDelay(500*time.Millisecond),
	)

	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		panic(err)
	}

	// 启动自动续约
	stopCh := make(chan struct{})
	go AutoRefresh(ctx, mutex, stopCh)
	time.Sleep(20 * time.Second) // 模拟耗时业务
	// 停止续约
	close(stopCh)

	if _, err := mutex.UnlockContext(ctx); err != nil {
		panic(err)
	}
}

func AutoRefresh(ctx context.Context, mutex *redsync.Mutex, stopCh <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second) // 每隔 5s 续约一次
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 续约，延长锁的过期时间
			if ok, err := mutex.ExtendContext(ctx); !ok || err != nil {
				slog.Error("Failed to extend mutex", "err", err, "status", ok)
			}
		case <-stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

/**

$ docker run -d \
  -p 6379:6379 \
  --name redis \
  redis:8.4.0 \
  redis-server --requirepass "password" --appendonly yes

*/
