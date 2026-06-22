package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	redis    *redis.Client
	tokenGen func() string
}

func NewAuthService(r *redis.Client) *AuthService {
	return &AuthService{
		redis:    r,
		tokenGen: generateToken,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	// 验证用户名密码（省略）

	token := s.tokenGen()
	key := fmt.Sprintf("session:%s", token)
	err := s.redis.Set(ctx, key, username, time.Hour).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) GetUsernameByToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("session:%s", token)
	username, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return username, nil
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "fallback-token"
	}
	return hex.EncodeToString(b)
}
