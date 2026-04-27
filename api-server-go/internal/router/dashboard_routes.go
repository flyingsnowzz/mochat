package router

import "github.com/gin-gonic/gin"

// registerCorpRoutes 注册企业相关路由
func (r *Router) registerCorpRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.corpHandler.Index)
	group.GET("/select", r.corpHandler.Select)
	group.GET("/show", r.corpHandler.Show)
	group.POST("/store", r.corpHandler.Store)
	group.PUT("/update", r.corpHandler.Update)
	group.PUT("/update/:id", r.corpHandler.Update)
	group.POST("/bind", r.corpHandler.Bind)
	group.POST("/bind/:id", r.corpHandler.Bind)
	group.GET("/weWorkCallback", r.corpHandler.WeWorkCallback)
	group.POST("/weWorkCallback", r.corpHandler.WeWorkCallback)
}

// registerUserRoutes 注册用户管理路由
func (r *Router) registerUserRoutes(group *gin.RouterGroup) {
	group.POST("/auth", r.userHandler.Auth)
	group.GET("/loginShow", r.userHandler.LoginShow)
	group.POST("/logout", r.userHandler.Logout)

	// 需要权限的路由
	permGroup := group.Group("")
	permGroup.Use(r.permission)
	{
		permGroup.GET("/index", r.userHandler.Index)
		permGroup.GET("/show/:id", r.userHandler.Show)
		permGroup.POST("/store", r.userHandler.Store)
		permGroup.PUT("/update/:id", r.userHandler.Update)
		permGroup.PUT("/passwordUpdate/:id", r.userHandler.PasswordUpdate)
		permGroup.PUT("/passwordReset/:id", r.userHandler.PasswordReset)
		permGroup.PUT("/statusUpdate/:id", r.userHandler.StatusUpdate)
	}
}

// registerMenuRoutes 注册菜单管理路由
func (r *Router) registerMenuRoutes(group *gin.RouterGroup) {
	group.Use(r.permission)
	group.GET("/index", r.menuHandler.Index)
	group.GET("/select", r.menuHandler.Select)
	group.GET("/show", r.menuHandler.Show)
	group.GET("/show/:id", r.menuHandler.Show)
	group.GET("/iconIndex", r.menuHandler.IconIndex)
	group.POST("/store", r.menuHandler.Store)
	group.PUT("/update", r.menuHandler.Update)
	group.PUT("/update/:id", r.menuHandler.Update)
	group.DELETE("/destroy", r.menuHandler.Destroy)
	group.DELETE("/destroy/:id", r.menuHandler.Destroy)
	group.PUT("/statusUpdate", r.menuHandler.StatusUpdate)
	group.PUT("/statusUpdate/:id", r.menuHandler.StatusUpdate)
}

// registerRoleRoutes 注册角色管理路由
func (r *Router) registerRoleRoutes(group *gin.RouterGroup) {
	group.Use(r.permission)
	group.GET("/index", r.roleHandler.Index)
	group.GET("/select", r.roleHandler.Select)
	group.GET("/show/:id", r.roleHandler.Show)
	group.POST("/store", r.roleHandler.Store)
	group.PUT("/update/:id", r.roleHandler.Update)
	group.DELETE("/destroy/:id", r.roleHandler.Destroy)
	group.PUT("/statusUpdate/:id", r.roleHandler.StatusUpdate)
	group.GET("/permissionShow/:id", r.roleHandler.PermissionShow)
	group.POST("/permissionStore/:id", r.roleHandler.PermissionStore)
	group.GET("/permissionByUser", r.roleHandler.PermissionByUser)
	group.GET("/showEmployee/:id", r.roleHandler.ShowEmployee)
}

// registerContactRoutes 注册联系人管理路由
func (r *Router) registerContactRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.contactHandler.Index)
	group.GET("/show", r.contactHandler.Show)
	group.GET("/show/:id", r.contactHandler.Show)
	group.PUT("/update", r.contactHandler.Update)
	group.PUT("/update/:id", r.contactHandler.Update)
	group.PUT("/synContact", r.contactHandler.SynContact)
	group.POST("/synContact", r.contactHandler.SynContact)
	group.GET("/track", r.contactHandler.Track)
	group.GET("/track/:id", r.contactHandler.Track)
	group.GET("/lossContact", r.contactHandler.LossContact)
	group.GET("/source", r.contactHandler.Source)
	group.POST("/batchLabeling", r.contactHandler.BatchLabeling)
}

// registerContactFieldRoutes 注册联系人字段管理路由
func (r *Router) registerContactFieldRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.contactFieldHandler.Index)
	group.GET("/show/:id", r.contactFieldHandler.Show)
	group.GET("/portrait", r.contactFieldHandler.Portrait)
	group.POST("/store", r.contactFieldHandler.Store)
	group.PUT("/update", r.contactFieldHandler.Update)
	group.PUT("/update/:id", r.contactFieldHandler.Update)
	group.PUT("/batchUpdate", r.contactFieldHandler.BatchUpdate)
	group.PUT("/statusUpdate", r.contactFieldHandler.StatusUpdate)
	group.PUT("/statusUpdate/:id", r.contactFieldHandler.StatusUpdate)
	group.DELETE("/destroy", r.contactFieldHandler.Destroy)
	group.DELETE("/destroy/:id", r.contactFieldHandler.Destroy)
	group.GET("/portrait/:id", r.contactFieldHandler.Portrait)
}

// registerContactFieldPivotRoutes 注册联系人字段关联路由
func (r *Router) registerContactFieldPivotRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.contactFieldPivotHandler.Index)
	group.PUT("/update/:id", r.contactFieldPivotHandler.Update)
}

// registerContactTagRoutes 注册联系人标签管理路由
func (r *Router) registerContactTagRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.tagHandler.Index)
	group.GET("/allTag", r.tagHandler.AllTag)
	group.GET("/detail/:id", r.tagHandler.Detail)
	group.GET("/contactTagList", r.tagHandler.ContactTagList)
	group.POST("/store", r.tagHandler.Store)
	group.PUT("/update/:id", r.tagHandler.Update)
	group.PUT("/move/:id", r.tagHandler.Move)
	group.DELETE("/destroy/:id", r.tagHandler.Destroy)
	group.POST("/synContactTag", r.tagHandler.SynContactTag)
}

// registerContactTagGroupRoutes 注册联系人标签组管理路由
func (r *Router) registerContactTagGroupRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.contactTagGroupHandler.Index)
	group.GET("/detail/:id", r.contactTagGroupHandler.Detail)
	group.POST("/store", r.contactTagGroupHandler.Store)
	group.PUT("/update/:id", r.contactTagGroupHandler.Update)
	group.DELETE("/destroy/:id", r.contactTagGroupHandler.Destroy)
}

// registerContactRoomRoutes 注册联系人所在群聊路由
func (r *Router) registerContactRoomRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.contactRoomHandler.Index)
}

// registerDeptRoutes 注册部门管理路由
func (r *Router) registerDeptRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.deptHandler.Index)
	group.GET("/pageIndex", r.deptHandler.PageIndex)
	group.GET("/selectByPhone", r.deptHandler.SelectByPhone)
	group.GET("/showEmployee", r.deptHandler.ShowEmployee)
	group.GET("/showEmployee/:id", r.deptHandler.ShowEmployee)
}

// registerEmployeeRoutes 注册员工管理路由
func (r *Router) registerEmployeeRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.employeeHandler.Index)
	group.GET("/searchCondition", r.employeeHandler.SearchCondition)
	group.POST("/syncEmployee", r.employeeHandler.SyncEmployee)
	group.PUT("/syncEmployee", r.employeeHandler.SyncEmployee)
	group.PUT("/synEmployee", r.employeeHandler.SyncEmployee)
	group.POST("/synEmployee", r.employeeHandler.SyncEmployee)
	group.GET("/department/memberIndex", r.deptHandler.DepartmentMemberIndex)
}

// registerRoomRoutes 注册客户群管理路由
func (r *Router) registerRoomRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.roomHandler.Index)
	group.GET("/roomIndex", r.roomHandler.RoomIndex)
	group.PUT("/batchUpdate", r.roomHandler.BatchUpdate)
	group.GET("/statistics", r.roomHandler.Statistics)
	group.GET("/statisticsIndex", r.roomHandler.StatisticsIndex)
	group.POST("/sync", r.roomHandler.Sync)
}

// registerRoomGroupRoutes 注册客户群分组管理路由
func (r *Router) registerRoomGroupRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.roomGroupHandler.Index)
	group.POST("/store", r.roomGroupHandler.Store)
	group.PUT("/update/:id", r.roomGroupHandler.Update)
	group.DELETE("/destroy/:id", r.roomGroupHandler.Destroy)
}

// registerAgentRoutes 注册第三方应用管理路由
func (r *Router) registerAgentRoutes(group *gin.RouterGroup) {
	group.POST("/store", r.agentHandler.Store)
	group.GET("/txtVerifyShow", r.agentHandler.TxtVerifyShow)
	group.POST("/txtVerifyUpload", r.agentHandler.TxtVerifyUpload)
}

// registerMediumRoutes 注册素材管理路由
func (r *Router) registerMediumRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.mediumHandler.Index)
	group.GET("/show/:id", r.mediumHandler.Show)
	group.GET("/show", r.mediumHandler.Show)
	group.POST("/store", r.mediumHandler.Store)
	group.PUT("/update/:id", r.mediumHandler.Update)
	group.PUT("/update", r.mediumHandler.Update)
	group.DELETE("/destroy/:id", r.mediumHandler.Destroy)
	group.DELETE("/destroy", r.mediumHandler.Destroy)
	group.PUT("/groupUpdate/:id", r.mediumHandler.GroupUpdate)
	group.PUT("/groupUpdate", r.mediumHandler.GroupUpdate)
}

// registerMediumGroupRoutes 注册素材分组管理路由
func (r *Router) registerMediumGroupRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.mediumGroupHandler.Index)
	group.POST("/store", r.mediumGroupHandler.Store)
	group.PUT("/update/:id", r.mediumGroupHandler.Update)
	group.DELETE("/destroy/:id", r.mediumGroupHandler.Destroy)
}

// registerIndexRoutes 注册首页数据路由
func (r *Router) registerIndexRoutes(group *gin.RouterGroup) {
	group.GET("/index/index", r.indexHandler.Index)
	group.GET("/index/lineChat", r.indexHandler.LineChat)
	group.GET("/corpData/index", r.indexHandler.Index)
	group.GET("/corpData/lineChat", r.indexHandler.LineChat)
}

// registerChatToolRoutes 注册聊天工具路由
func (r *Router) registerChatToolRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.chatToolHandler.Index)
}

// registerOfficialAccountRoutes 注册公众号管理路由
func (r *Router) registerOfficialAccountRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.officialAccountHandler.Index)
	group.GET("/getPreAuthUrl", r.officialAccountHandler.GetPreAuthUrl)
	group.GET("/authEventCallback", r.officialAccountHandler.AuthEventCallback)
	group.POST("/authEventCallback", r.officialAccountHandler.AuthEventCallback)
	group.POST("/set", r.officialAccountHandler.Set)
}

// registerCommonRoutes 注册公共接口路由
func (r *Router) registerCommonRoutes(group *gin.RouterGroup) {
	group.POST("/upload", r.commonHandler.Upload)
	group.POST("/uploadFile", r.commonHandler.UploadFile)
}

// registerGreetingRoutes 注册欢迎语管理路由
func (r *Router) registerGreetingRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.greetingHandler.Index)
	group.GET("/show/:id", r.greetingHandler.Show)
	group.GET("/show", r.greetingHandler.Show)
	group.POST("/store", r.greetingHandler.Store)
	group.PUT("/update/:id", r.greetingHandler.Update)
	group.PUT("/update", r.greetingHandler.Update)
	group.DELETE("/destroy/:id", r.greetingHandler.Destroy)
	group.DELETE("/destroy", r.greetingHandler.Destroy)
}

// registerChannelCodeRoutes 注册渠道码管理路由
func (r *Router) registerChannelCodeRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.channelCodeHandler.Index)
	group.GET("/contact", r.channelCodeHandler.Contact)
	group.GET("/statistics", r.channelCodeHandler.Statistics)
	group.GET("/statisticsIndex", r.channelCodeHandler.StatisticsIndex)
	group.GET("/show", r.channelCodeHandler.Show)
	group.GET("/show/:id", r.channelCodeHandler.Show)
	group.POST("/store", r.channelCodeHandler.Store)
	group.PUT("/update", r.channelCodeHandler.Update)
	group.PUT("/update/:id", r.channelCodeHandler.Update)
}

// registerChannelCodeGroupRoutes 注册渠道码分组管理路由
func (r *Router) registerChannelCodeGroupRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.channelCodeHandler.GroupIndex)
	group.POST("/store", r.channelCodeHandler.GroupStore)
	group.PUT("/update", r.channelCodeHandler.GroupUpdate)
	group.PUT("/move", r.channelCodeHandler.GroupMove)
}

// registerStatisticRoutes 注册统计数据路由
func (r *Router) registerStatisticRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.statisticHandler.Index)
	group.GET("/employees", r.statisticHandler.Employees)
	group.GET("/topList", r.statisticHandler.TopList)
	group.GET("/employeesTrend", r.statisticHandler.EmployeesTrend)
	group.GET("/employeeCounts", r.statisticHandler.EmployeeCounts)
}

// registerRoomWelcomeRoutes 注册进群欢迎语路由
func (r *Router) registerRoomWelcomeRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.roomWelcomeHandler.Index)
	group.GET("/show", r.roomWelcomeHandler.Show)
	group.GET("/show/:id", r.roomWelcomeHandler.Show)
	group.POST("/store", r.roomWelcomeHandler.Store)
	group.PUT("/update", r.roomWelcomeHandler.Update)
	group.PUT("/update/:id", r.roomWelcomeHandler.Update)
	group.DELETE("/destroy", r.roomWelcomeHandler.Destroy)
	group.DELETE("/destroy/:id", r.roomWelcomeHandler.Destroy)
}

// registerWorkRoomAutoPullRoutes 注册自动拉群路由
func (r *Router) registerWorkRoomAutoPullRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.workRoomAutoPullHandler.Index)
	group.GET("/show", r.workRoomAutoPullHandler.Show)
	group.GET("/show/:id", r.workRoomAutoPullHandler.Show)
	group.POST("/store", r.workRoomAutoPullHandler.Store)
	group.PUT("/update", r.workRoomAutoPullHandler.Update)
	group.PUT("/update/:id", r.workRoomAutoPullHandler.Update)
	group.DELETE("/destroy", r.workRoomAutoPullHandler.Destroy)
	group.DELETE("/destroy/:id", r.workRoomAutoPullHandler.Destroy)
}

// registerRoomTagPullRoutes 注册群标签拉取路由
func (r *Router) registerRoomTagPullRoutes(group *gin.RouterGroup) {
	group.GET("/index", r.roomTagPullHandler.Index)
	group.POST("/create", r.roomTagPullHandler.Create)
	group.GET("/detail", r.roomTagPullHandler.Detail)
	group.GET("/contactDetail", r.roomTagPullHandler.ContactDetail)
}
