package sidebar

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

// WorkRoomHandler 侧边栏群聊处理器
// 处理侧边栏相关的群聊操作，包括群聊管理等

type WorkRoomHandler struct {
	db      *gorm.DB                // 数据库连接
	roomSvc *service.WorkRoomService // 群聊服务实例
}

// NewWorkRoomHandler 创建群聊处理器实例
// 参数:
//   - db: 数据库连接
// 返回值:
//   - *WorkRoomHandler: 群聊处理器实例

func NewWorkRoomHandler(db *gorm.DB) *WorkRoomHandler {
	return &WorkRoomHandler{
		db:      db,
		roomSvc: service.NewWorkRoomService(db),
	}
}

// RoomManage 群聊管理
// 请求方法: GET
// 请求路径: /sidebar/workRoom/roomManage
// 响应:
//   - 成功: 包含群聊列表的对象，每个群聊包含成员数量
//   - 失败: 错误信息

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
