package plugin

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

type GreetingHandler struct{}

func NewGreetingHandler() *GreetingHandler {
	return &GreetingHandler{}
}

func (h *GreetingHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *GreetingHandler) Show(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *GreetingHandler) Store(c *gin.Context) {
	response.SuccessMsg(c, "创建成功")
}

func (h *GreetingHandler) Update(c *gin.Context) {
	response.SuccessMsg(c, "更新成功")
}

func (h *GreetingHandler) Destroy(c *gin.Context) {
	response.SuccessMsg(c, "删除成功")
}