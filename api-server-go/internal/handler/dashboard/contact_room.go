// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含群聊成员管理的处理器：
// WorkContactRoomHandler - 处理群聊成员的查询操作
package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

// WorkContactRoomHandler 群聊成员管理处理器
// 提供群聊成员的查询功能
// 主要职责：
// 1. 获取群聊的成员列表
//
// 依赖服务：
// - WorkContactRoomService: 群聊成员服务

type WorkContactRoomHandler struct {
	roomContactSvc *service.WorkContactRoomService // 群聊成员服务
}

// NewWorkContactRoomHandler 创建群聊成员管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：群聊成员管理处理器实例
func NewWorkContactRoomHandler(db *gorm.DB) *WorkContactRoomHandler {
	return &WorkContactRoomHandler{
		roomContactSvc: service.NewWorkContactRoomService(db),
	}
}

// Index 获取群聊成员列表
// 根据群聊 ID 获取其成员列表，支持分页
// 处理流程：
// 1. 获取 roomId 查询参数
// 2. 解析 roomId 为 uint 类型
// 3. 获取分页参数
// 4. 调用服务层获取群聊成员列表
// 5. 返回分页结果
// 参数：
//
//	roomId - 群聊 ID（查询参数）
//	page - 页码，默认为 1
//	pageSize - 每页数量，默认为 20
//
// 返回：包含群聊成员列表、总数、分页信息的响应
func (h *WorkContactRoomHandler) Index(c *gin.Context) {
	// 获取 roomId 查询参数
	roomIDStr := c.Query("roomId")
	if roomIDStr == "" {
		response.Fail(c, response.ErrParams, "缺少 roomId 参数")
		return
	}

	// 解析 roomId 为 uint 类型
	roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
	if err != nil {
		response.Fail(c, response.ErrParams, "roomId 参数错误")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 调用服务层获取群聊成员列表
	contacts, total, err := h.roomContactSvc.List(uint(roomID), page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取群聊成员列表失败")
		return
	}

	// 返回分页结果
	response.PageResult(c, contacts, total, page, pageSize)
}
