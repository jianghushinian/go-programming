package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:generate mockgen -source store.go -destination gomock.go -package main

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserStore interface {
	Create(u *User) error
	Get(id int) (*User, error)
}

func NewDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db}
}

type userStore struct {
	db *gorm.DB
}

func (s *userStore) Create(u *User) error {
	return s.db.Save(u).Error
}

func (s *userStore) Get(id int) (u *User, err error) {
	return u, s.db.First(&u, id).Error
}
