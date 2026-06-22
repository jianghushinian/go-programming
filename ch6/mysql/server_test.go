package main

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type fakeStore struct {
	users map[int]*User
}

func (f *fakeStore) Get(id int) (*User, error) {
	if u, ok := f.users[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("fake not found")
}

func (f *fakeStore) Create(u *User) error { return nil }

func TestUserHandler_GetUser_Fake(t *testing.T) {
	r := gin.Default()
	handler := &UserHandler{store: &fakeStore{
		users: map[int]*User{1: {ID: 1, Name: "江湖十年"}},
	}}
	handler.InitRoutes(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"id":1,"name":"江湖十年"}`, w.Body.String())
}

func TestUserHandler_GetUser_Mock(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockUserStore := NewMockUserStore(ctrl)
	mockUserStore.EXPECT().Get(2).Return(&User{
		ID:   2,
		Name: "jianghushinian",
	}, nil).Times(1)

	r := gin.Default()
	handler := &UserHandler{store: mockUserStore}
	handler.InitRoutes(r)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"id":2,"name":"jianghushinian"}`, w.Body.String())
}
