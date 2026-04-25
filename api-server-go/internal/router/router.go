package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/config"
	"mochat-api-server/internal/handler/dashboard"
	dashboardPlugin "mochat-api-server/internal/handler/dashboard/plugin"
	"mochat-api-server/internal/handler/sidebar"
	"mochat-api-server/internal/middleware"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/pkg/storage"
	pluginService "mochat-api-server/internal/service/plugin"
)

func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CoreMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLogMiddleware())

	r.Static("/storage", "./storage/upload")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	dashboardWhiteRoutes := []string{
		"/dashboard/user/auth",
		"/dashboard/user/loginShow",
		"/dashboard/corp/weWorkCallback",
		"/dashboard/officialAccount/authEventCallback",
		"/dashboard/*/officialAccount/messageEventCallback",
		"/dashboard/common/upload",
		"/dashboard/common/uploadFile",
		"/operation/*",
		"/health",
	}

	sidebarWhiteRoutes := []string{
		"/sidebar/wxJSSDK/config",
		"/sidebar/common/upload",
		"/sidebar/agent/auth",
		"/sidebar/agent/oauth",
		"/sidebar/agent/jssdkConfig",
		"/health",
	}

	dashboardAuth := middleware.DashboardAuthMiddleware(cfg.JWT, dashboardWhiteRoutes)
	sidebarAuth := middleware.SidebarAuthMiddleware(cfg.JWT, sidebarWhiteRoutes)
	permission := middleware.PermissionMiddleware(db)

	corpHandler := dashboard.NewCorpHandler(db)
	userHandler := dashboard.NewUserHandler(db, cfg.JWT)
	menuHandler := dashboard.NewMenuHandler(db)
	roleHandler := dashboard.NewRoleHandler(db)
	contactHandler := dashboard.NewContactHandler(db)
	tagHandler := dashboard.NewContactTagHandler(db)
	contactFieldHandler := dashboard.NewContactFieldHandler(db)
	contactFieldPivotHandler := dashboard.NewContactFieldPivotHandler(db)
	contactTagGroupHandler := dashboard.NewWorkContactTagGroupHandler(db)
	contactRoomHandler := dashboard.NewWorkContactRoomHandler(db)
	deptHandler := dashboard.NewDepartmentHandler(db)
	employeeHandler := dashboard.NewEmployeeHandler(db)
	roomHandler := dashboard.NewRoomHandler(db)
	roomGroupHandler := dashboard.NewRoomGroupHandler(db)
	agentHandler := dashboard.NewAgentHandler(db)
	mediumHandler := dashboard.NewMediumHandler(db)
	mediumGroupHandler := dashboard.NewMediumGroupHandler(db)
	indexHandler := dashboard.NewIndexHandler(db)
	chatToolHandler := dashboard.NewChatToolHandler(db)
	officialAccountHandler := dashboard.NewOfficialAccountHandler()
	greetingHandler := dashboard.NewGreetingHandler(db)
	contactMessageBatchSendHandler := dashboardPlugin.NewContactMessageBatchSendHandler(
		db,
		pluginService.NewContactMessageBatchSendService(db),
	)
	channelCodeHandler := dashboardPlugin.NewChannelCodeHandler(
		db,
		pluginService.NewChannelCodeService(db),
		pluginService.NewChannelCodeGroupService(db),
	)
	statisticHandler := dashboardPlugin.NewStatisticHandler()
	workFissionHandler := dashboardPlugin.NewWorkFissionHandler(db, pluginService.NewWorkFissionService(db))

	commonHandler := dashboard.NewCommonHandler()

	sidebarContactHandler := sidebar.NewWorkContactHandler(db)
	sidebarRoomHandler := sidebar.NewWorkRoomHandler(db)
	sidebarAgentHandler := sidebar.NewWorkAgentHandler()
	sidebarMediumHandler := sidebar.NewMediumHandler()
	sidebarMediumGroupHandler := sidebar.NewMediumGroupHandler()
	sidebarCommonHandler := sidebar.NewCommonHandler(storage.DefaultStorage)

	dashboardGroup := r.Group("/dashboard")
	dashboardGroup.Use(dashboardAuth)
	{
		dashboardGroup.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

		corpGroup := dashboardGroup.Group("/corp")
		{
			corpGroup.GET("/index", corpHandler.Index)
			corpGroup.GET("/select", corpHandler.Select)
			corpGroup.GET("/show", corpHandler.Show)
			corpGroup.GET("/show/:id", corpHandler.Show)
			corpGroup.POST("/store", corpHandler.Store)
			corpGroup.PUT("/update", corpHandler.Update)
			corpGroup.PUT("/update/:id", corpHandler.Update)
			corpGroup.POST("/bind", corpHandler.Bind)
			corpGroup.POST("/bind/:id", corpHandler.Bind)
			corpGroup.GET("/weWorkCallback", corpHandler.WeWorkCallback)
			corpGroup.POST("/weWorkCallback", corpHandler.WeWorkCallback)
		}

		userGroup := dashboardGroup.Group("/user")
		{
			userGroup.POST("/auth", userHandler.Auth)
			userGroup.GET("/loginShow", userHandler.LoginShow)
			userGroup.POST("/logout", userHandler.Logout)
			userGroup.Use(permission)
			userGroup.GET("/index", userHandler.Index)
			userGroup.GET("/show/:id", userHandler.Show)
			userGroup.POST("/store", userHandler.Store)
			userGroup.PUT("/update/:id", userHandler.Update)
			userGroup.PUT("/passwordUpdate/:id", userHandler.PasswordUpdate)
			userGroup.PUT("/passwordReset/:id", userHandler.PasswordReset)
			userGroup.PUT("/statusUpdate/:id", userHandler.StatusUpdate)
		}

		menuGroup := dashboardGroup.Group("/menu")
		menuGroup.Use(permission)
		{
			menuGroup.GET("/index", menuHandler.Index)
			menuGroup.GET("/select", menuHandler.Select)
			menuGroup.GET("/show", menuHandler.Show)
			menuGroup.GET("/show/:id", menuHandler.Show)
			menuGroup.GET("/iconIndex", menuHandler.IconIndex)
			menuGroup.POST("/store", menuHandler.Store)
			menuGroup.PUT("/update", menuHandler.Update)
			menuGroup.PUT("/update/:id", menuHandler.Update)
			menuGroup.DELETE("/destroy", menuHandler.Destroy)
			menuGroup.DELETE("/destroy/:id", menuHandler.Destroy)
			menuGroup.PUT("/statusUpdate", menuHandler.StatusUpdate)
			menuGroup.PUT("/statusUpdate/:id", menuHandler.StatusUpdate)
		}

		roleGroup := dashboardGroup.Group("/role")
		roleGroup.Use(permission)
		{
			roleGroup.GET("/index", roleHandler.Index)
			roleGroup.GET("/select", roleHandler.Select)
			roleGroup.GET("/show/:id", roleHandler.Show)
			roleGroup.POST("/store", roleHandler.Store)
			roleGroup.PUT("/update/:id", roleHandler.Update)
			roleGroup.DELETE("/destroy/:id", roleHandler.Destroy)
			roleGroup.PUT("/statusUpdate/:id", roleHandler.StatusUpdate)
			roleGroup.GET("/permissionShow/:id", roleHandler.PermissionShow)
			roleGroup.POST("/permissionStore/:id", roleHandler.PermissionStore)
			roleGroup.GET("/permissionByUser", roleHandler.PermissionByUser)
			roleGroup.GET("/showEmployee/:id", roleHandler.ShowEmployee)
		}

		contactGroup := dashboardGroup.Group("/workContact")
		{
			contactGroup.GET("/index", contactHandler.Index)
			contactGroup.GET("/show", contactHandler.Show)
			contactGroup.GET("/show/:id", contactHandler.Show)
			contactGroup.PUT("/update", contactHandler.Update)
			contactGroup.PUT("/update/:id", contactHandler.Update)
			contactGroup.PUT("/synContact", contactHandler.SynContact)
			contactGroup.POST("/synContact", contactHandler.SynContact)
			contactGroup.GET("/track", contactHandler.Track)
			contactGroup.GET("/track/:id", contactHandler.Track)
			contactGroup.GET("/lossContact", contactHandler.LossContact)
			contactGroup.GET("/source", contactHandler.Source)
			contactGroup.POST("/batchLabeling", contactHandler.BatchLabeling)
		}

		contactFieldGroup := dashboardGroup.Group("/contactField")
		{
			contactFieldGroup.GET("/index", contactFieldHandler.Index)
			contactFieldGroup.GET("/show/:id", contactFieldHandler.Show)
			contactFieldGroup.GET("/portrait", contactFieldHandler.Portrait)
			contactFieldGroup.POST("/store", contactFieldHandler.Store)
			contactFieldGroup.PUT("/update", contactFieldHandler.Update)
			contactFieldGroup.PUT("/update/:id", contactFieldHandler.Update)
			contactFieldGroup.PUT("/batchUpdate", contactFieldHandler.BatchUpdate)
			contactFieldGroup.PUT("/statusUpdate", contactFieldHandler.StatusUpdate)
			contactFieldGroup.PUT("/statusUpdate/:id", contactFieldHandler.StatusUpdate)
			contactFieldGroup.DELETE("/destroy", contactFieldHandler.Destroy)
			contactFieldGroup.DELETE("/destroy/:id", contactFieldHandler.Destroy)
			contactFieldGroup.GET("/portrait/:id", contactFieldHandler.Portrait)
		}

		contactFieldPivotGroup := dashboardGroup.Group("/contactFieldPivot")
		{
			contactFieldPivotGroup.GET("/index", contactFieldPivotHandler.Index)
			contactFieldPivotGroup.PUT("/update/:id", contactFieldPivotHandler.Update)
		}

		contactTagGroup := dashboardGroup.Group("/workContactTag")
		{
			contactTagGroup.GET("/index", tagHandler.Index)
			contactTagGroup.GET("/allTag", tagHandler.AllTag)
			contactTagGroup.GET("/detail/:id", tagHandler.Detail)
			contactTagGroup.GET("/contactTagList", tagHandler.ContactTagList)
			contactTagGroup.POST("/store", tagHandler.Store)
			contactTagGroup.PUT("/update/:id", tagHandler.Update)
			contactTagGroup.PUT("/move/:id", tagHandler.Move)
			contactTagGroup.DELETE("/destroy/:id", tagHandler.Destroy)
			contactTagGroup.POST("/synContactTag", tagHandler.SynContactTag)
		}

		contactTagGroupGroup := dashboardGroup.Group("/workContactTagGroup")
		{
			contactTagGroupGroup.GET("/index", contactTagGroupHandler.Index)
			contactTagGroupGroup.GET("/detail/:id", contactTagGroupHandler.Detail)
			contactTagGroupGroup.POST("/store", contactTagGroupHandler.Store)
			contactTagGroupGroup.PUT("/update/:id", contactTagGroupHandler.Update)
			contactTagGroupGroup.DELETE("/destroy/:id", contactTagGroupHandler.Destroy)
		}

		contactRoomGroup := dashboardGroup.Group("/workContactRoom")
		{
			contactRoomGroup.GET("/index", contactRoomHandler.Index)
		}

		deptGroup := dashboardGroup.Group("/workDepartment")
		{
			deptGroup.GET("/index", deptHandler.Index)
			deptGroup.GET("/pageIndex", deptHandler.PageIndex)
			deptGroup.GET("/selectByPhone", deptHandler.SelectByPhone)
			deptGroup.GET("/showEmployee", deptHandler.ShowEmployee)
			deptGroup.GET("/showEmployee/:id", deptHandler.ShowEmployee)
		}

		employeeGroup := dashboardGroup.Group("/workEmployee")
		{
			employeeGroup.GET("/index", employeeHandler.Index)
			employeeGroup.GET("/searchCondition", employeeHandler.SearchCondition)
			employeeGroup.POST("/syncEmployee", employeeHandler.SyncEmployee)
			employeeGroup.PUT("/syncEmployee", employeeHandler.SyncEmployee)
			employeeGroup.PUT("/synEmployee", employeeHandler.SyncEmployee)
			employeeGroup.POST("/synEmployee", employeeHandler.SyncEmployee)
			employeeGroup.GET("/department/memberIndex", deptHandler.DepartmentMemberIndex)
		}

		roomGroup := dashboardGroup.Group("/workRoom")
		{
			roomGroup.GET("/index", roomHandler.Index)
			roomGroup.GET("/roomIndex", roomHandler.RoomIndex)
			roomGroup.PUT("/batchUpdate", roomHandler.BatchUpdate)
			roomGroup.GET("/statistics", roomHandler.Statistics)
			roomGroup.GET("/statisticsIndex", roomHandler.StatisticsIndex)
			roomGroup.POST("/sync", roomHandler.Sync)
		}

		roomGroupSubGroup := dashboardGroup.Group("/workRoomGroup")
		{
			roomGroupSubGroup.GET("/index", roomGroupHandler.Index)
			roomGroupSubGroup.POST("/store", roomGroupHandler.Store)
			roomGroupSubGroup.PUT("/update/:id", roomGroupHandler.Update)
			roomGroupSubGroup.DELETE("/destroy/:id", roomGroupHandler.Destroy)
		}

		agentGroup := dashboardGroup.Group("/agent")
		{
			agentGroup.POST("/store", agentHandler.Store)
			agentGroup.GET("/txtVerifyShow", agentHandler.TxtVerifyShow)
			agentGroup.POST("/txtVerifyUpload", agentHandler.TxtVerifyUpload)
		}

		mediumGroup := dashboardGroup.Group("/medium")
		{
			mediumGroup.GET("/index", mediumHandler.Index)
			mediumGroup.GET("/show/:id", mediumHandler.Show)
			mediumGroup.GET("/show", mediumHandler.Show)
			mediumGroup.POST("/store", mediumHandler.Store)
			mediumGroup.PUT("/update/:id", mediumHandler.Update)
			mediumGroup.PUT("/update", mediumHandler.Update)
			mediumGroup.DELETE("/destroy/:id", mediumHandler.Destroy)
			mediumGroup.DELETE("/destroy", mediumHandler.Destroy)
			mediumGroup.PUT("/groupUpdate/:id", mediumHandler.GroupUpdate)
			mediumGroup.PUT("/groupUpdate", mediumHandler.GroupUpdate)
		}

		mediumGroupSubGroup := dashboardGroup.Group("/mediumGroup")
		{
			mediumGroupSubGroup.GET("/index", mediumGroupHandler.Index)
			mediumGroupSubGroup.POST("/store", mediumGroupHandler.Store)
			mediumGroupSubGroup.PUT("/update/:id", mediumGroupHandler.Update)
			mediumGroupSubGroup.DELETE("/destroy/:id", mediumGroupHandler.Destroy)
		}

		dashboardGroup.GET("/index/index", indexHandler.Index)
		dashboardGroup.GET("/index/lineChat", indexHandler.LineChat)
		dashboardGroup.GET("/corpData/index", indexHandler.Index)
		dashboardGroup.GET("/corpData/lineChat", indexHandler.LineChat)

		chatToolGroup := dashboardGroup.Group("/chatTool")
		{
			chatToolGroup.GET("/index", chatToolHandler.Index)
		}

		oaGroup := dashboardGroup.Group("/officialAccount")
		{
			oaGroup.GET("/index", officialAccountHandler.Index)
			oaGroup.GET("/getPreAuthUrl", officialAccountHandler.GetPreAuthUrl)
			oaGroup.GET("/authEventCallback", officialAccountHandler.AuthEventCallback)
			oaGroup.POST("/authEventCallback", officialAccountHandler.AuthEventCallback)
			oaGroup.POST("/set", officialAccountHandler.Set)
		}

		commonGroup := dashboardGroup.Group("/common")
		{
			commonGroup.POST("/upload", commonHandler.Upload)
			commonGroup.POST("/uploadFile", commonHandler.UploadFile)
		}

		greetingGroup := dashboardGroup.Group("/greeting")
		{
			greetingGroup.GET("/index", greetingHandler.Index)
			greetingGroup.GET("/show/:id", greetingHandler.Show)
			greetingGroup.GET("/show", greetingHandler.Show)
			greetingGroup.POST("/store", greetingHandler.Store)
			greetingGroup.PUT("/update/:id", greetingHandler.Update)
			greetingGroup.PUT("/update", greetingHandler.Update)
			greetingGroup.DELETE("/destroy/:id", greetingHandler.Destroy)
			greetingGroup.DELETE("/destroy", greetingHandler.Destroy)
		}

		contactMessageBatchSendGroup := dashboardGroup.Group("/contactMessageBatchSend")
		{
			contactMessageBatchSendGroup.GET("/index", contactMessageBatchSendHandler.Index)
			contactMessageBatchSendGroup.POST("/store", contactMessageBatchSendHandler.Store)
			contactMessageBatchSendGroup.GET("/show", contactMessageBatchSendHandler.Show)
			contactMessageBatchSendGroup.GET("/show/:id", contactMessageBatchSendHandler.Show)
			contactMessageBatchSendGroup.DELETE("/destroy", contactMessageBatchSendHandler.Destroy)
			contactMessageBatchSendGroup.DELETE("/destroy/:id", contactMessageBatchSendHandler.Destroy)
			contactMessageBatchSendGroup.GET("/messageShow", contactMessageBatchSendHandler.MessageShow)
			contactMessageBatchSendGroup.POST("/remind", contactMessageBatchSendHandler.Remind)
			contactMessageBatchSendGroup.GET("/contactReceiveIndex", contactMessageBatchSendHandler.ContactReceiveIndex)
			contactMessageBatchSendGroup.GET("/employeeSendIndex", contactMessageBatchSendHandler.EmployeeSendIndex)
		}

		channelCodeGroup := dashboardGroup.Group("/channelCode")
		{
			channelCodeGroup.GET("/index", channelCodeHandler.Index)
			channelCodeGroup.GET("/contact", channelCodeHandler.Contact)
			channelCodeGroup.GET("/statistics", channelCodeHandler.Statistics)
			channelCodeGroup.GET("/statisticsIndex", channelCodeHandler.StatisticsIndex)
			channelCodeGroup.GET("/show", channelCodeHandler.Show)
			channelCodeGroup.GET("/show/:id", channelCodeHandler.Show)
			channelCodeGroup.POST("/store", channelCodeHandler.Store)
			channelCodeGroup.PUT("/update", channelCodeHandler.Update)
			channelCodeGroup.PUT("/update/:id", channelCodeHandler.Update)
		}

		channelCodeGroupGroup := dashboardGroup.Group("/channelCodeGroup")
		{
			channelCodeGroupGroup.GET("/index", channelCodeHandler.GroupIndex)
			channelCodeGroupGroup.POST("/store", channelCodeHandler.GroupStore)
			channelCodeGroupGroup.PUT("/update", channelCodeHandler.GroupUpdate)
			channelCodeGroupGroup.PUT("/move", channelCodeHandler.GroupMove)
		}

		statisticGroup := dashboardGroup.Group("/statistic")
		{
			statisticGroup.GET("/index", statisticHandler.Index)
			statisticGroup.GET("/employees", statisticHandler.Employees)
			statisticGroup.GET("/topList", statisticHandler.TopList)
		}

		fissionGroup := dashboardGroup.Group("/workFission")
		{
			fissionGroup.GET("/index", workFissionHandler.Index)
			fissionGroup.GET("/statistics", workFissionHandler.Statistics)
			fissionGroup.GET("/inviteData", workFissionHandler.InviteData)
			fissionGroup.GET("/inviteDetail", workFissionHandler.InviteDetail)
			fissionGroup.GET("/chooseContact", workFissionHandler.ChooseContact)
			fissionGroup.POST("/invite", workFissionHandler.Invite)
			fissionGroup.GET("/show", workFissionHandler.Show)
			fissionGroup.GET("/show/:id", workFissionHandler.Show)
			fissionGroup.GET("/info", workFissionHandler.Info)
			fissionGroup.GET("/info/:id", workFissionHandler.Info)
			fissionGroup.POST("/store", workFissionHandler.Store)
			fissionGroup.PUT("/update", workFissionHandler.Update)
			fissionGroup.PUT("/update/:id", workFissionHandler.Update)
		}
	}

	sidebarGroup := r.Group("/sidebar")
	sidebarGroup.Use(sidebarAuth)
	{
		sidebarGroup.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
		sidebarGroup.GET("/workContact/show/:id", sidebarContactHandler.Show)
		sidebarGroup.GET("/workContact/detail/:id", sidebarContactHandler.Detail)
		sidebarGroup.PUT("/workContact/update/:id", sidebarContactHandler.Update)
		sidebarGroup.GET("/workContact/track/:id", sidebarContactHandler.Track)
		sidebarGroup.GET("/workRoom/roomManage", sidebarRoomHandler.RoomManage)
		sidebarGroup.GET("/agent/auth", sidebarAgentHandler.Auth)
		sidebarGroup.GET("/agent/oauth", sidebarAgentHandler.OAuth)
		sidebarGroup.GET("/agent/jssdkConfig", sidebarAgentHandler.JssdkConfig)
		sidebarGroup.GET("/medium/index", sidebarMediumHandler.Index)
		sidebarGroup.PUT("/medium/mediaIdUpdate/:id", sidebarMediumHandler.MediaIdUpdate)
		sidebarGroup.GET("/mediumGroup/index", sidebarMediumGroupHandler.Index)
		sidebarGroup.POST("/common/upload", sidebarCommonHandler.Upload)
		sidebarGroup.GET("/wxJSSDK/config", sidebarCommonHandler.WxJSSDKConfig)
	}

	operationGroup := r.Group("/operation")
	{
		operationGroup.GET("/officialAccount/authRedirect", func(c *gin.Context) {
			response.Success(c, gin.H{"redirect": "/"})
		})
	}

	return r
}
