package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type RoomHandler struct {
	svc      *service.WorkRoomService
	groupSvc *service.WorkRoomGroupService
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{
		svc:      service.NewWorkRoomService(db),
		groupSvc: service.NewWorkRoomGroupService(db),
	}
}

func (h *RoomHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	rooms, total, err := h.svc.List(corpID.(uint), page, pageSize, nil)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取客户群列表失败")
		return
	}
	response.PageResult(c, rooms, total, page, pageSize)
}

func (h *RoomHandler) RoomIndex(c *gin.Context) { h.Index(c) }

func (h *RoomHandler) BatchUpdate(c *gin.Context) {
	var req struct {
		IDs     []uint `json:"ids"`
		GroupID uint   `json:"groupId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.BatchUpdateGroup(req.IDs, req.GroupID); err != nil {
		response.Fail(c, response.ErrDB, "批量修改失败")
		return
	}
	response.SuccessMsg(c, "批量修改成功")
}

func (h *RoomHandler) Statistics(c *gin.Context) {
	response.Success(c, gin.H{})
}

func (h *RoomHandler) StatisticsIndex(c *gin.Context) {
	response.Success(c, []interface{}{})
}

func (h *RoomHandler) Sync(c *gin.Context) {
	response.SuccessMsg(c, "同步任务已提交")
}

type RoomGroupHandler struct {
	svc *service.WorkRoomGroupService
}

func NewRoomGroupHandler(db *gorm.DB) *RoomGroupHandler {
	return &RoomGroupHandler{svc: service.NewWorkRoomGroupService(db)}
}

func (h *RoomGroupHandler) Index(c *gin.Context) {
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	groups, err := h.svc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取群分组列表失败")
		return
	}
	response.Success(c, groups)
}

func (h *RoomGroupHandler) Store(c *gin.Context) {
	var group model.WorkRoomGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&group); err != nil {
		response.Fail(c, response.ErrDB, "创建群分组失败")
		return
	}
	response.Success(c, group)
}

func (h *RoomGroupHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	group := &model.WorkRoomGroup{}
	group.ID = uint(id)
	if err := c.ShouldBindJSON(group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Update(group); err != nil {
		response.Fail(c, response.ErrDB, "更新群分组失败")
		return
	}
	response.Success(c, group)
}

func (h *RoomGroupHandler) Destroy(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除群分组失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}
