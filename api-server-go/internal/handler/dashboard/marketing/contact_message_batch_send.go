package marketing

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type ContactMessageBatchSendHandler struct {
	service *plugin.ContactMessageBatchSendService
}

func NewContactMessageBatchSendHandler(service *plugin.ContactMessageBatchSendService) *ContactMessageBatchSendHandler {
	return &ContactMessageBatchSendHandler{service: service}
}

func (h *ContactMessageBatchSendHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *ContactMessageBatchSendHandler) Store(c *gin.Context) {
	response.SuccessMsg(c, "创建成功")
}

func (h *ContactMessageBatchSendHandler) Show(c *gin.Context) {
	response.Success(c, nil)
}

func (h *ContactMessageBatchSendHandler) Destroy(c *gin.Context) {
	response.SuccessMsg(c, "删除成功")
}

func (h *ContactMessageBatchSendHandler) Remind(c *gin.Context) {
	response.SuccessMsg(c, "提醒成功")
}

func (h *ContactMessageBatchSendHandler) ShowRoom(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *ContactMessageBatchSendHandler) EmployeeSendIndex(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *ContactMessageBatchSendHandler) ContactReceiveIndex(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}
