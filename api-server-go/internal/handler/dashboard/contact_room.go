package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type WorkContactRoomHandler struct {
	roomContactSvc *service.WorkContactRoomService
}

func NewWorkContactRoomHandler(db *gorm.DB) *WorkContactRoomHandler {
	return &WorkContactRoomHandler{
		roomContactSvc: service.NewWorkContactRoomService(db),
	}
}

func (h *WorkContactRoomHandler) Index(c *gin.Context) {
	roomIDStr := c.Query("roomId")
	if roomIDStr == "" {
		response.Fail(c, response.ErrParams, "缺少 roomId 参数")
		return
	}

	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		response.Fail(c, response.ErrParams, "roomId 参数错误")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	contacts, total, err := h.roomContactSvc.List(uint(roomID), page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取群聊成员列表失败")
		return
	}

	response.PageResult(c, contacts, total, page, pageSize)
}
