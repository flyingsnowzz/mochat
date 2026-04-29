package marketing

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type ContactTransferHandler struct {
	service *plugin.WorkTransferService
}

func NewContactTransferHandler(service *plugin.WorkTransferService) *ContactTransferHandler {
	return &ContactTransferHandler{service: service}
}

func (h *ContactTransferHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *ContactTransferHandler) Info(c *gin.Context) {
	response.Success(c, nil)
}

func (h *ContactTransferHandler) Log(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *ContactTransferHandler) Room(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *ContactTransferHandler) UnassignedList(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *ContactTransferHandler) SaveUnassignedList(c *gin.Context) {
	response.SuccessMsg(c, "保存成功")
}

func (h *ContactTransferHandler) TransferRoom(c *gin.Context) {
	response.SuccessMsg(c, "转移成功")
}
