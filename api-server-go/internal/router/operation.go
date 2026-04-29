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

	// 任务宝运营端路由
	r.registerWorkFissionOperationRoutes(group.Group("/workFission"))
}

func (r *Router) registerWorkFissionOperationRoutes(group *gin.RouterGroup) {
	group.GET("/auth", r.workFissionHandler.Auth)
	group.POST("/inviteFriends", r.workFissionHandler.InviteFriends)
	group.POST("/taskData", r.workFissionHandler.TaskData)
	group.POST("/receive", r.workFissionHandler.Receive)
	group.GET("/poster", r.workFissionHandler.Poster)
	group.GET("/openUserInfo", r.workFissionHandler.OpenUserInfo)
}
