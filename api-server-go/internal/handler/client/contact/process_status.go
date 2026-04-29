package contact

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
)

type ProcessStatusHandler struct{}

func NewProcessStatusHandler() *ProcessStatusHandler {
	return &ProcessStatusHandler{}
}

func (h *ProcessStatusHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *ProcessStatusHandler) Update(c *gin.Context) {
	response.SuccessMsg(c, "更新成功")
}
