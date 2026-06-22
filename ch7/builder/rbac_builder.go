package builder

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Role 角色
type Role string

// Permission 权限
type Permission string

// RBAC 记录角色到权限的映射
type RBAC struct {
	grants map[Role]map[Permission]struct{}
}

// RBACBuilder 用于构建 RBAC 配置
type RBACBuilder struct {
	RBAC
}

// NewRBACBuilder 创建 Builder
func NewRBACBuilder() *RBACBuilder {
	return &RBACBuilder{RBAC{grants: make(map[Role]map[Permission]struct{})}}
}

// Allow 为角色添加权限
func (b *RBACBuilder) Allow(role Role, perm Permission) *RBACBuilder {
	if _, ok := b.grants[role]; !ok {
		b.grants[role] = make(map[Permission]struct{})
	}
	b.grants[role][perm] = struct{}{}
	return b
}

// Build 返回 RBAC 实例
func (b *RBACBuilder) Build() *RBAC {
	return &b.RBAC
}

// Has 判断角色是否包含权限
func (r *RBAC) Has(role Role, perm Permission) bool {
	if ps, ok := r.grants[role]; ok {
		if _, ok := ps[perm]; ok {
			return true
		}
	}
	return false
}

// Require 返回 Gin 中间件，检查请求角色是否拥有所需权限
// 此示例从请求头 X-User-Role 读取角色，实际项目可从上下文 / JWT / session 等获取
func (r *RBAC) Require(perm Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := Role(c.GetHeader("X-User-Role"))
		if role == "" || !r.Has(role, perm) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
