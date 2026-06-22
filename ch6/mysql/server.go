package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	store UserStore
}

func (h *UserHandler) InitRoutes(r *gin.Engine) {
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.store.Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func main() {
	db, err := NewDB("user:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}
	store := NewUserStore(db)
	handler := &UserHandler{store: store}
	r := gin.Default()
	handler.InitRoutes(r)
	r.Run(":8080")
}
