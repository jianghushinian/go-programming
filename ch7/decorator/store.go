//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID   int
	Name string
}

type UserStore interface {
	Create(u *User) error
	Get(id int) (*User, error)
}

type LoggingStore struct {
	s UserStore
}

func (l *LoggingStore) Create(u *User) error {
	log.Printf("Creating user: %+v", u)
	return l.s.Create(u)
}

func (l *LoggingStore) Get(id int) (*User, error) {
	log.Printf("Getting user by id: %d", id)
	return l.s.Get(id)
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

func main() {
	db, err := NewDB("user:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}
	store := NewUserStore(db)
	loggingStore := &LoggingStore{s: store}
	_ = loggingStore.Create(&User{ID: 1, Name: "江湖十年"})
	_, _ = loggingStore.Get(1)
}
