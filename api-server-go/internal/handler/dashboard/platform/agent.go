package platform

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/pkg/response"
)

// AgentHandler 企业微信应用处理器
// 处理企业微信应用相关操作，包括配置、授权等
type AgentHandler struct {
	db *gorm.DB
}

// NewAgentHandler 创建企业微信应用处理器实例
func NewAgentHandler(db *gorm.DB) *AgentHandler {
	return &AgentHandler{db: db}
}

// Index 获取企业微信应用列表
func (h *AgentHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

// Get 获取企业微信应用详情
func (h *AgentHandler) Get(c *gin.Context) {
	response.Success(c, gin.H{})
}

// Create 创建企业微信应用
func (h *AgentHandler) Create(c *gin.Context) {
	response.SuccessMsg(c, "创建成功")
}

// Update 更新企业微信应用
func (h *AgentHandler) Update(c *gin.Context) {
	response.SuccessMsg(c, "更新成功")
}

// Delete 删除企业微信应用
func (h *AgentHandler) Delete(c *gin.Context) {
	response.SuccessMsg(c, "删除成功")
}

// Store 存储企业微信应用
func (h *AgentHandler) Store(c *gin.Context) {
	response.SuccessMsg(c, "存储成功")
}

// TxtVerifyShow 显示文本验证
func (h *AgentHandler) TxtVerifyShow(c *gin.Context) {
	response.Success(c, gin.H{})
}

// TxtVerifyUpload 上传文本验证
func (h *AgentHandler) TxtVerifyUpload(c *gin.Context) {
	response.SuccessMsg(c, "上传成功")
}

// GetAuthUrl 获取授权URL
func (h *AgentHandler) GetAuthUrl(c *gin.Context) {
	response.Success(c, gin.H{
		"url": "https://open.work.weixin.qq.com/wwopen/sso/qrConnect",
	})
}

// AuthEventCallback 授权事件回调
func (h *AgentHandler) AuthEventCallback(c *gin.Context) {
	c.String(200, "success")
}
