package router

import "github.com/gin-gonic/gin"

// RegisterSidebarRoutes 注册 Sidebar 相关路由
func (r *Router) RegisterSidebarRoutes(group *gin.RouterGroup) {
	group.GET("/health", r.HealthCheck)

	// 侧边栏联系人相关路由
	group.GET("/workContact/show/:id", r.sidebarContactHandler.Show)
	group.GET("/workContact/detail/:id", r.sidebarContactHandler.Detail)
	group.PUT("/workContact/update/:id", r.sidebarContactHandler.Update)
	group.GET("/workContact/track/:id", r.sidebarContactHandler.Track)

	// 侧边栏客户群相关路由
	group.GET("/workRoom/roomManage", r.sidebarRoomHandler.RoomManage)

	// 侧边栏第三方应用路由
	group.GET("/agent/auth", r.sidebarAgentHandler.Auth)
	group.GET("/agent/oauth", r.sidebarAgentHandler.OAuth)
	group.GET("/agent/jssdkConfig", r.sidebarAgentHandler.JssdkConfig)

	// 侧边栏素材相关路由
	group.GET("/medium/index", r.sidebarMediumHandler.Index)
	group.PUT("/medium/mediaIdUpdate/:id", r.sidebarMediumHandler.MediaIdUpdate)
	group.GET("/mediumGroup/index", r.sidebarMediumGroupHandler.Index)

	// 侧边栏公共接口路由
	group.POST("/common/upload", r.sidebarCommonHandler.Upload)
	group.GET("/wxJSSDK/config", r.sidebarCommonHandler.WxJSSDKConfig)
}
