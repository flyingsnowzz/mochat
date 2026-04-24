package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
)

const (
	// DataPermissionAll 全部数据权限
	DataPermissionAll = 1
	// DataPermissionDept 部门数据权限
	DataPermissionDept = 2
	// DataPermissionSelf 仅本人数据权限
	DataPermissionSelf = 3
)

// PermissionMiddleware RBAC 权限中间件
func PermissionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		guard, _ := c.Get("guard")
		if guard == "sidebar" {
			c.Next()
			return
		}

		userID, exists := c.Get("userId")
		if !exists {
			c.Next()
			return
		}

		uid := userID.(uint)

		var user model.User
		if err := db.Select("id", "isSuperAdmin").First(&user, uid).Error; err != nil {
			response.FailWithHTTP(c, http.StatusForbidden, response.ErrPermission, "无权限")
			c.Abort()
			return
		}

		if user.IsSuperAdmin == 1 {
			c.Set("dataPermission", 0)
			c.Set("roleIds", []uint{})
			c.Next()
			return
		}

		var userRoles []model.RbacUserRole
		if err := db.Where("user_id = ?", uid).Find(&userRoles).Error; err != nil {
			response.FailWithHTTP(c, http.StatusForbidden, response.ErrPermission, "无权限")
			c.Abort()
			return
		}

		if len(userRoles) == 0 {
			response.FailWithHTTP(c, http.StatusForbidden, response.ErrPermission, "无权限")
			c.Abort()
			return
		}

		roleIDs := make([]uint, len(userRoles))
		for i, ur := range userRoles {
			roleIDs[i] = ur.RoleID
		}

		var role model.RbacRole
		if err := db.Where("id IN ? AND status = 1", roleIDs).First(&role).Error; err != nil {
			response.FailWithHTTP(c, http.StatusForbidden, response.ErrPermission, "无权限")
			c.Abort()
			return
		}

		c.Set("dataPermission", role.DataPermission)
		c.Set("roleIds", roleIDs)
		c.Next()
	}
}

// GetDataPermission 从上下文中获取数据权限级别
func GetDataPermission(c *gin.Context) int {
	if val, exists := c.Get("dataPermission"); exists {
		return val.(int)
	}
	return DataPermissionSelf
}
