package sidebar

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type WorkRoomHandler struct {
	db      *gorm.DB
	roomSvc *service.WorkRoomService
}

func NewWorkRoomHandler(db *gorm.DB) *WorkRoomHandler {
	return &WorkRoomHandler{
		db:      db,
		roomSvc: service.NewWorkRoomService(db),
	}
}

func (h *WorkRoomHandler) RoomManage(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		response.Fail(c, response.ErrParams, "未获取到企业信息")
		return
	}

	// 获取群聊列表
	var rooms []model.WorkRoom
	if err := h.db.Where("corp_id = ?", corpID).Find(&rooms).Error; err != nil {
		response.Fail(c, response.ErrDB, "获取群聊列表失败")
		return
	}

	// 获取每个群聊的成员数量
	for i := range rooms {
		var memberCount int64
		h.db.Model(&model.WorkContactRoom{}).Where("room_id = ?", rooms[i].ID).Count(&memberCount)
		rooms[i].RoomMax = int(memberCount)
	}

	response.Success(c, gin.H{
		"rooms": rooms,
	})
}
