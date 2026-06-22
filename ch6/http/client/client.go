package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// User 用户信息结构体
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserClient 用户客户端
type UserClient struct {
	BaseURL string
}

// NewUserClient 创建用户客户端实例
func NewUserClient(baseURL string) *UserClient {
	return &UserClient{BaseURL: baseURL}
}

// GetUser 从 HTTP 服务获取用户信息
func (c *UserClient) GetUser(userID int) (*User, error) {
	url := fmt.Sprintf("%s/users/%d", c.BaseURL, userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API status code: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func main() {
	client := &UserClient{BaseURL: "http://localhost:8080"}
	user, err := client.GetUser(1)
	if err != nil {
		fmt.Printf("GetUser error: %v\n", err)
		return
	}
	fmt.Printf("User: ID=%d, Name=%s\n", user.ID, user.Name)
}
