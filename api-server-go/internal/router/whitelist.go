package router

// getDashboardWhiteRoutes 获取 Dashboard 白名单路由
func getDashboardWhiteRoutes() []string {
	return []string{
		"/dashboard/user/auth",
		"/dashboard/user/loginShow",
		"/dashboard/corp/weWorkCallback",
		"/dashboard/officialAccount/authEventCallback",
		"/dashboard/*/officialAccount/messageEventCallback",
		"/operation/*",
		"/health",
	}
}

// getSidebarWhiteRoutes 获取 Sidebar 白名单路由
func getSidebarWhiteRoutes() []string {
	return []string{
		"/sidebar/wxJSSDK/config",
		"/sidebar/common/upload",
		"/sidebar/agent/auth",
		"/sidebar/agent/oauth",
		"/sidebar/agent/jssdkConfig",
		"/health",
	}
}
