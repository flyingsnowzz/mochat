package router

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

// RegisterOperationRoutes 注册运营操作路由
func (r *Router) RegisterOperationRoutes(group *gin.RouterGroup) {
	group.GET("/officialAccount/authRedirect", func(c *gin.Context) {
		response.Success(c, gin.H{"redirect": "/"})
	})
}
