package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/config"
	dashboardSystem "mochat-api-server/internal/handler/dashboard/system"
	dashboardContact "mochat-api-server/internal/handler/dashboard/contact"
	dashboardOrg "mochat-api-server/internal/handler/dashboard/organization"
	dashboardContent "mochat-api-server/internal/handler/dashboard/content"
	dashboardMarketing "mochat-api-server/internal/handler/dashboard/marketing"
	dashboardAnalysis "mochat-api-server/internal/handler/dashboard/analysis"
	dashboardPlatform "mochat-api-server/internal/handler/dashboard/platform"
	dashboardCommon "mochat-api-server/internal/handler/dashboard/common"
	clientContact "mochat-api-server/internal/handler/client/contact"
	clientOrg "mochat-api-server/internal/handler/client/organization"
	clientContent "mochat-api-server/internal/handler/client/content"
	clientPlatform "mochat-api-server/internal/handler/client/platform"
	clientCommon "mochat-api-server/internal/handler/client/common"
	"mochat-api-server/internal/middleware"
)

// Router 路由器结构体，持有所有依赖和路由配置
type Router struct {
	engine          *gin.Engine
	config          *config.Config
	db              *gorm.DB
	dashboardAuth   gin.HandlerFunc
	sidebarAuth     gin.HandlerFunc
	permission      gin.HandlerFunc

	// Dashboard System Handlers (系统管理)
	corpHandler             *dashboardSystem.CorpHandler
	userHandler             *dashboardSystem.UserHandler
	menuHandler             *dashboardSystem.MenuHandler
	roleHandler             *dashboardSystem.RoleHandler

	// Dashboard Contact Handlers (客户管理)
	contactHandler          *dashboardContact.ContactHandler
	tagHandler              *dashboardContact.ContactTagHandler
	contactFieldHandler     *dashboardContact.ContactFieldHandler
	contactFieldPivotHandler *dashboardContact.ContactFieldPivotHandler
	contactTagGroupHandler  *dashboardContact.WorkContactTagGroupHandler
	contactRoomHandler      *dashboardContact.WorkContactRoomHandler

	// Dashboard Organization Handlers (组织架构)
	deptHandler             *dashboardOrg.DepartmentHandler
	employeeHandler         *dashboardOrg.EmployeeHandler
	roomHandler             *dashboardOrg.RoomHandler
	roomGroupHandler        *dashboardOrg.RoomGroupHandler

	// Dashboard Content Handlers (内容管理)
	mediumHandler           *dashboardContent.MediumHandler
	mediumGroupHandler      *dashboardContent.MediumGroupHandler
	greetingHandler         *dashboardContent.GreetingHandler

	// Dashboard Analysis Handlers (数据分析)
	indexHandler            *dashboardAnalysis.DashboardIndexHandler

	// Dashboard Platform Handlers (平台配置)
	agentHandler            *dashboardPlatform.AgentHandler
	officialAccountHandler  *dashboardPlatform.OfficialAccountHandler

	// Dashboard Common Handlers
	commonHandler    *dashboardCommon.CommonHandler
	chatToolHandler  *dashboardCommon.ChatToolHandler

	// Marketing Handlers
	channelCodeHandler      *dashboardMarketing.ChannelCodeHandler
	statisticHandler        *dashboardAnalysis.StatisticHandler
	roomWelcomeHandler      *dashboardMarketing.RoomWelcomeHandler
	workRoomAutoPullHandler *dashboardMarketing.WorkRoomAutoPullHandler
	roomTagPullHandler      *dashboardMarketing.RoomTagPullHandler

	// Client Handlers
	sidebarContactHandler    *clientContact.WorkContactHandler
	sidebarRoomHandler       *clientOrg.WorkRoomHandler
	sidebarAgentHandler      *clientPlatform.WorkAgentHandler
	sidebarMediumHandler     *clientContent.MediumHandler
	sidebarMediumGroupHandler *clientContent.MediumGroupHandler
	sidebarCommonHandler     *clientCommon.CommonHandler

	storage                any // 存储接口，具体类型根据实际使用确定
}

// NewRouter 创建新的路由器实例
func NewRouter(cfg *config.Config, db *gorm.DB) *Router {
	return &Router{
		engine:          gin.New(),
		config:          cfg,
		db:              db,
		dashboardAuth:   middleware.DashboardAuthMiddleware(cfg.JWT, getDashboardWhiteRoutes()),
		sidebarAuth:     middleware.SidebarAuthMiddleware(cfg.JWT, getSidebarWhiteRoutes()),
		permission:      middleware.PermissionMiddleware(db),
	}
}

// Setup 初始化路由器
func (r *Router) Setup() *gin.Engine {
	r.setMode()
	r.useMiddlewares()
	r.registerStaticRoutes()
	r.registerHealthRoutes()
	r.initHandlers()
	r.registerDashboardRoutes()
	r.registerSidebarRoutes()
	r.registerOperationRoutes()
	return r.engine
}

// setMode 设置 Gin 运行模式
func (r *Router) setMode() {
	if r.config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// useMiddlewares 使用全局中间件
func (r *Router) useMiddlewares() {
	r.engine.Use(gin.Recovery())
	r.engine.Use(middleware.CoreMiddleware())
	r.engine.Use(middleware.CORSMiddleware())
	r.engine.Use(middleware.RequestLogMiddleware())
}

// registerStaticRoutes 注册静态文件路由
func (r *Router) registerStaticRoutes() {
	r.engine.Static("/storage", "./storage/upload")
}

// registerHealthRoutes 注册健康检查路由
func (r *Router) registerHealthRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// registerDashboardRoutes 注册 Dashboard 相关路由
func (r *Router) registerDashboardRoutes() {
	dashboardGroup := r.engine.Group("/dashboard")
	dashboardGroup.Use(r.dashboardAuth)
	r.RegisterDashboardRoutes(dashboardGroup)
}

// registerSidebarRoutes 注册 Sidebar 相关路由
func (r *Router) registerSidebarRoutes() {
	sidebarGroup := r.engine.Group("/sidebar")
	sidebarGroup.Use(r.sidebarAuth)
	r.RegisterSidebarRoutes(sidebarGroup)
}

// registerOperationRoutes 注册运营操作路由
func (r *Router) registerOperationRoutes() {
	operationGroup := r.engine.Group("/operation")
	r.RegisterOperationRoutes(operationGroup)
}
