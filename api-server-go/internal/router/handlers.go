package router

import (
	clientCommon "mochat-api-server/internal/handler/client/common"
	clientContact "mochat-api-server/internal/handler/client/contact"
	clientContent "mochat-api-server/internal/handler/client/content"
	clientOrg "mochat-api-server/internal/handler/client/organization"
	clientPlatform "mochat-api-server/internal/handler/client/platform"
	dashboardAnalysis "mochat-api-server/internal/handler/dashboard/analysis"
	dashboardCommon "mochat-api-server/internal/handler/dashboard/common"
	dashboardContact "mochat-api-server/internal/handler/dashboard/contact"
	dashboardContent "mochat-api-server/internal/handler/dashboard/content"
	dashboardMarketing "mochat-api-server/internal/handler/dashboard/marketing"
	dashboardOrg "mochat-api-server/internal/handler/dashboard/organization"
	dashboardPlatform "mochat-api-server/internal/handler/dashboard/platform"
	dashboardSystem "mochat-api-server/internal/handler/dashboard/system"
	businessService "mochat-api-server/internal/service/business"
	pluginService "mochat-api-server/internal/service/plugin"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/pkg/storage"
)

// initHandlers 初始化所有处理器
func (r *Router) initHandlers() {
	// ============ System 模块 (系统管理) ============
	corpSvc := businessService.NewCorpService(r.db)
	userSvc := businessService.NewUserService(r.db)
	jwtCfg := response.JWTConfig{
		DashboardSecret: r.config.JWT.DashboardSecret,
		DashboardPrefix: r.config.JWT.DashboardPrefix,
		SidebarSecret:   r.config.JWT.SidebarSecret,
	}
	r.corpHandler = dashboardSystem.NewCorpHandler(corpSvc)
	r.userHandler = dashboardSystem.NewUserHandler(userSvc, corpSvc, jwtCfg)
	r.menuHandler = dashboardSystem.NewMenuHandler(r.db)
	r.roleHandler = dashboardSystem.NewRoleHandler(r.db)

	// ============ Contact 模块 (客户管理) ============
	r.contactHandler = dashboardContact.NewContactHandler(r.db)
	r.tagHandler = dashboardContact.NewContactTagHandler(r.db)
	r.contactFieldHandler = dashboardContact.NewContactFieldHandler(r.db)
	r.contactFieldPivotHandler = dashboardContact.NewContactFieldPivotHandler(r.db)
	r.contactTagGroupHandler = dashboardContact.NewWorkContactTagGroupHandler(r.db)
	r.contactRoomHandler = dashboardContact.NewWorkContactRoomHandler(r.db)

	// ============ Organization 模块 (组织架构) ============
	r.deptHandler = dashboardOrg.NewDepartmentHandler(r.db)
	r.employeeHandler = dashboardOrg.NewEmployeeHandler(r.db)
	r.roomHandler = dashboardOrg.NewRoomHandler(r.db)
	r.roomGroupHandler = dashboardOrg.NewRoomGroupHandler(r.db)

	// ============ Content 模块 (内容管理) ============
	r.mediumHandler = dashboardContent.NewMediumHandler(r.db)
	r.mediumGroupHandler = dashboardContent.NewMediumGroupHandler(r.db)
	r.greetingHandler = dashboardContent.NewGreetingHandler(r.db)

	// ============ Marketing 模块 (营销工具) ============
	r.channelCodeHandler = dashboardMarketing.NewChannelCodeHandler(
		r.db,
		pluginService.NewChannelCodeService(r.db),
		pluginService.NewChannelCodeGroupService(r.db),
	)
	r.roomWelcomeHandler = dashboardMarketing.NewRoomWelcomeHandler(pluginService.NewRoomWelcomeService(r.db))
	r.workRoomAutoPullHandler = dashboardMarketing.NewWorkRoomAutoPullHandler(pluginService.NewWorkRoomAutoPullService(r.db))
	r.roomTagPullHandler = dashboardMarketing.NewRoomTagPullHandler(pluginService.NewRoomTagPullService(r.db))

	// ============ Analysis 模块 (数据分析) ============
	r.indexHandler = dashboardAnalysis.NewDashboardIndexHandler(r.db)
	r.statisticHandler = dashboardAnalysis.NewStatisticHandler()

	// ============ Platform 模块 (平台配置) ============
	r.agentHandler = dashboardPlatform.NewAgentHandler(r.db)
	r.officialAccountHandler = dashboardPlatform.NewOfficialAccountHandler()

	// ============ Common 公共接口 ============
	r.commonHandler = dashboardCommon.NewCommonHandler()
	r.chatToolHandler = dashboardCommon.NewChatToolHandler(r.db)

	// ============ Sidebar/Client 模块 (客户端) ============
	r.sidebarContactHandler = clientContact.NewWorkContactHandler(r.db)
	r.sidebarRoomHandler = clientOrg.NewWorkRoomHandler(r.db)
	r.sidebarAgentHandler = clientPlatform.NewWorkAgentHandler()
	r.sidebarMediumHandler = clientContent.NewMediumHandler()
	r.sidebarMediumGroupHandler = clientContent.NewMediumGroupHandler()
	r.sidebarCommonHandler = clientCommon.NewCommonHandler(storage.DefaultStorage)
}
