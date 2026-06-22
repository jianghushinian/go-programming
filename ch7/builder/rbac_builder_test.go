package builder

import (
	"log"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRBACBuilder(t *testing.T) {
	rbac := NewRBACBuilder().
		Allow("admin", "create").
		Allow("admin", "read").
		Allow("admin", "update").
		Allow("guest", "read").
		Build()

	r := gin.Default()
	r.POST("/users", rbac.Require("create"), CreateUserHandler)
	r.GET("/users/:id", rbac.Require("read"), GetUserHandler)

	// NOTE: Run forever.
	r.Run(":8080")
}

func CreateUserHandler(c *gin.Context) {
	log.Println("CreateUserHandler")
}

func GetUserHandler(c *gin.Context) {
	log.Println("GetUserHandler")
}
