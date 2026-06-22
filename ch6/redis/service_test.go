package main

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Login_GetUsername(t *testing.T) {
	ctx := context.Background()
	mockClient, mock := redismock.NewClientMock()
	svc := &AuthService{
		redis:    mockClient,
		tokenGen: func() string { return "tok123" },
	}

	mock.ExpectSet("session:tok123", "江湖十年", time.Hour).SetVal("OK")
	mock.ExpectGet("session:tok123").SetVal("江湖十年")

	token, err := svc.Login(ctx, "江湖十年", "password")
	assert.NoError(t, err)
	assert.Equal(t, "tok123", token)

	username, err := svc.GetUsernameByToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, "江湖十年", username)
}
