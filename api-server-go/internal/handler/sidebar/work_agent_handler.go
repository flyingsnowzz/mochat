package sidebar

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

// WorkAgentHandler 侧边栏企业微信应用处理器
// 处理侧边栏相关的企业微信应用操作，包括认证、OAuth和JSSDK配置等

type WorkAgentHandler struct{}

// NewWorkAgentHandler 创建企业微信应用处理器实例
// 返回值:
//   - *WorkAgentHandler: 企业微信应用处理器实例

func NewWorkAgentHandler() *WorkAgentHandler {
	return &WorkAgentHandler{}
}

// Auth 企业微信应用认证
// 请求方法: GET
// 请求路径: /sidebar/agent/auth
// 响应:
//   - 成功: 空对象

func (h *WorkAgentHandler) Auth(c *gin.Context) {
	response.Success(c, gin.H{})
}

// OAuth 企业微信应用OAuth
// 请求方法: GET
// 请求路径: /sidebar/agent/oauth
// 响应:
//   - 成功: 包含重定向URL的对象

func (h *WorkAgentHandler) OAuth(c *gin.Context) {
	response.Success(c, gin.H{"redirect": "/"})
}

// JssdkConfig 企业微信应用JSSDK配置
// 请求方法: GET
// 请求路径: /sidebar/agent/jssdkConfig
// 响应:
//   - 成功: 包含JSSDK配置的对象

func (h *WorkAgentHandler) JssdkConfig(c *gin.Context) {
	response.Success(c, gin.H{
		"appId":     "wx1234567890",
		"timestamp": 1234567890,
		"nonceStr":  "randomstring",
		"signature": "abcdefg",
	})
}
