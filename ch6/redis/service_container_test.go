package main

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func RunRedisInContainer(t *testing.T) (*redis.Client, func()) {
	ctx := context.Background()
	// 启动 Redis 容器
	redisC, err := testcontainers.Run(
		ctx, "redis:6.0.20-alpine", // 指定容器镜像
		testcontainers.WithExposedPorts("6379/tcp"), // 指定容器暴露端口
		testcontainers.WithWaitStrategy( // 等待容器内 Redis 就绪
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		),
	)
	assert.NoError(t, err)
	// 获取容器中 Redis 连接地址
	endpoint, _ := redisC.Endpoint(ctx, "")
	client := redis.NewClient(&redis.Options{Addr: endpoint})
	return client, func() {
		testcontainers.CleanupContainer(t, redisC) // 清理容器
	}
}

func TestRunWithRedisInContainer(t *testing.T) {
	client, cleanup := RunRedisInContainer(t)
	defer cleanup()

	ctx := context.Background()
	// svc := &AuthService{
	// 	redis:    client,
	// 	tokenGen: generateToken,
	// }
	svc := NewAuthService(client)

	token, err := svc.Login(ctx, "江湖十年", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, token) // token 不为空

	username, err := svc.GetUsernameByToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, "江湖十年", username)
}
