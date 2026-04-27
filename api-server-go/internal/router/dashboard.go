package router

import "github.com/gin-gonic/gin"

// HealthCheck 健康检查接口
func (r *Router) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

// RegisterDashboardRoutes 注册 Dashboard 相关路由
func (r *Router) RegisterDashboardRoutes(group *gin.RouterGroup) {
	// 健康检查
	group.GET("/health", r.HealthCheck)

	// 企业相关路由
	r.registerCorpRoutes(group.Group("/corp"))

	// 用户管理路由
	r.registerUserRoutes(group.Group("/user"))

	// 菜单管理路由
	r.registerMenuRoutes(group.Group("/menu"))

	// 角色管理路由
	r.registerRoleRoutes(group.Group("/role"))

	// 联系人管理路由
	r.registerContactRoutes(group.Group("/workContact"))

	// 联系人字段管理路由
	r.registerContactFieldRoutes(group.Group("/contactField"))

	// 联系人字段关联路由
	r.registerContactFieldPivotRoutes(group.Group("/contactFieldPivot"))

	// 联系人标签管理路由
	r.registerContactTagRoutes(group.Group("/workContactTag"))

	// 联系人标签组管理路由
	r.registerContactTagGroupRoutes(group.Group("/workContactTagGroup"))

	// 联系人所在群聊路由
	r.registerContactRoomRoutes(group.Group("/workContactRoom"))

	// 部门管理路由
	r.registerDeptRoutes(group.Group("/workDepartment"))

	// 员工管理路由
	r.registerEmployeeRoutes(group.Group("/workEmployee"))

	// 客户群管理路由
	r.registerRoomRoutes(group.Group("/workRoom"))

	// 客户群分组管理路由
	r.registerRoomGroupRoutes(group.Group("/workRoomGroup"))

	// 第三方应用管理路由
	r.registerAgentRoutes(group.Group("/agent"))

	// 素材管理路由
	r.registerMediumRoutes(group.Group("/medium"))

	// 素材分组管理路由
	r.registerMediumGroupRoutes(group.Group("/mediumGroup"))

	// 首页数据路由
	r.registerIndexRoutes(group)

	// 聊天工具路由
	r.registerChatToolRoutes(group.Group("/chatTool"))

	// 公众号管理路由
	r.registerOfficialAccountRoutes(group.Group("/officialAccount"))

	// 公共接口路由
	r.registerCommonRoutes(group.Group("/common"))

	// 欢迎语管理路由
	r.registerGreetingRoutes(group.Group("/greeting"))

	// 渠道码管理路由
	r.registerChannelCodeRoutes(group.Group("/channelCode"))

	// 渠道码分组管理路由
	r.registerChannelCodeGroupRoutes(group.Group("/channelCodeGroup"))

	// 统计数据路由
	r.registerStatisticRoutes(group.Group("/statistic"))

	// 进群欢迎语路由
	r.registerRoomWelcomeRoutes(group.Group("/roomWelcome"))

	// 自动拉群路由
	r.registerWorkRoomAutoPullRoutes(group.Group("/workRoomAutoPull"))

	// 群标签拉取路由
	r.registerRoomTagPullRoutes(group.Group("/roomTagPull"))
}
