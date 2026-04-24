package sidebar

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

type WorkAgentHandler struct{}

func NewWorkAgentHandler() *WorkAgentHandler {
	return &WorkAgentHandler{}
}

func (h *WorkAgentHandler) Auth(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *WorkAgentHandler) OAuth(c *gin.Context) {
	response.Success(c, gin.H{"redirect": "/"})
}

func (h *WorkAgentHandler) JssdkConfig(c *gin.Context) {
	response.Success(c, gin.H{
		"appId":     "wx1234567890",
		"timestamp": 1234567890,
		"nonceStr":  "randomstring",
		"signature": "abcdefg",
	})
}
