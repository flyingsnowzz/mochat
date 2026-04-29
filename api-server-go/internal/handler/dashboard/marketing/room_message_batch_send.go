package marketing

import (
	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type RoomMessageBatchSendHandler struct {
	service *plugin.RoomMessageBatchSendService
}

func NewRoomMessageBatchSendHandler(service *plugin.RoomMessageBatchSendService) *RoomMessageBatchSendHandler {
	return &RoomMessageBatchSendHandler{service: service}
}

func (h *RoomMessageBatchSendHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *RoomMessageBatchSendHandler) Store(c *gin.Context) {
	response.SuccessMsg(c, "创建成功")
}

func (h *RoomMessageBatchSendHandler) Show(c *gin.Context) {
	response.Success(c, nil)
}

func (h *RoomMessageBatchSendHandler) Destroy(c *gin.Context) {
	response.SuccessMsg(c, "删除成功")
}

func (h *RoomMessageBatchSendHandler) Remind(c *gin.Context) {
	response.SuccessMsg(c, "提醒成功")
}

func (h *RoomMessageBatchSendHandler) RoomOwnerSendIndex(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

func (h *RoomMessageBatchSendHandler) RoomReceiveIndex(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}
