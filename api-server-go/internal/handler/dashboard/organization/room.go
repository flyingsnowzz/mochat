package organization

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

type RoomHandler struct {
	db       *gorm.DB
	svc      *service.WorkRoomService
	groupSvc *service.WorkRoomGroupService
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{
		db:       db,
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

	// 转换为前端期望的格式
	type RoomResponse struct {
		WorkRoomID   uint   `json:"workRoomId"`
		RoomName     string `json:"roomName"`
		MemberNum    int64  `json:"memberNum"`
		OwnerName    string `json:"ownerName"`
		RoomGroup    string `json:"roomGroup"`
		StatusText   string `json:"statusText"`
		InRoomNum    int64  `json:"inRoomNum"`
		OutRoomNum   int64  `json:"outRoomNum"`
		Notice       string `json:"notice"`
		CreateTime   string `json:"createTime"`
	}

	var roomResponses []RoomResponse
	for _, room := range rooms {
		// 查询群成员数
		var memberNum int64
		h.db.Model(&model.WorkContactRoom{}).Where("room_id = ? AND status = 1", room.ID).Count(&memberNum)

		// 查询群主名称
		ownerName := ""
		var employee model.WorkEmployee
		h.db.Where("id = ?", room.OwnerID).First(&employee)
		if employee.ID > 0 {
			ownerName = employee.Name
		}

		// 查询分组名称
		roomGroup := "未分组"
		if room.RoomGroupID > 0 {
			var group model.WorkRoomGroup
			h.db.Where("id = ?", room.RoomGroupID).First(&group)
			if group.ID > 0 {
				roomGroup = group.Name
			}
		}

		// 状态文本
		statusText := "正常"
		switch room.Status {
		case 1:
			statusText = "跟进人离职"
		case 2:
			statusText = "离职继承中"
		case 3:
			statusText = "离职继承完成"
		}

		// 今日入群/退群数
		today := time.Now().Format("2006-01-02")
		var inRoomNum, outRoomNum int64
		h.db.Model(&model.WorkContactRoom{}).Where("room_id = ? AND DATE(join_time) = ?", room.ID, today).Count(&inRoomNum)
		h.db.Model(&model.WorkContactRoom{}).Where("room_id = ? AND DATE(out_time) = ?", room.ID, today).Count(&outRoomNum)

		roomResponses = append(roomResponses, RoomResponse{
			WorkRoomID: room.ID,
			RoomName:   room.Name,
			MemberNum:  memberNum,
			OwnerName:  ownerName,
			RoomGroup:  roomGroup,
			StatusText: statusText,
			InRoomNum:  inRoomNum,
			OutRoomNum: outRoomNum,
			Notice:     room.Notice,
			CreateTime: room.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	response.PageResult(c, roomResponses, total, page, pageSize)
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
	corpIDStr := c.Query("corpId")
	var corpID uint
	if corpIDStr != "" {
		parsedID, err := strconv.ParseUint(corpIDStr, 10, 32)
		if err == nil {
			corpID = uint(parsedID)
		}
	}
	if corpID == 0 {
		if cid, ok := c.Get("corpId"); ok {
			corpID = cid.(uint)
		}
	}
	groups, err := h.svc.List(corpID)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取群分组列表失败")
		return
	}

	// 转换为前端期望的格式
	type GroupResponse struct {
		WorkRoomGroupId   uint   `json:"workRoomGroupId"`
		WorkRoomGroupName string `json:"workRoomGroupName"`
	}

	var groupResponses []GroupResponse
	for _, group := range groups {
		groupResponses = append(groupResponses, GroupResponse{
			WorkRoomGroupId:   group.ID,
			WorkRoomGroupName: group.Name,
		})
	}

	response.Success(c, groupResponses)
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
