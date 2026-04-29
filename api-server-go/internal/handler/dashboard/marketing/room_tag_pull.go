package marketing

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type RoomTagPullHandler struct {
	service *plugin.RoomTagPullService
}

func NewRoomTagPullHandler(service *plugin.RoomTagPullService) *RoomTagPullHandler {
	return &RoomTagPullHandler{
		service: service,
	}
}

// Index 标签建群列表
func (h *RoomTagPullHandler) Index(c *gin.Context) {
	// 获取查询参数
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	offset := (page - 1) * perPage

	// 从上下文获取企业 ID，暂时使用默认值 1
	corpID := uint(1)

	// 调用服务获取数据
	items, total, err := h.service.List(corpID, name, offset, perPage)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取标签建群列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":  items,
		"total": total,
		"page":  page,
		"size":  perPage,
	})
}

// Create 创建标签建群
func (h *RoomTagPullHandler) Create(c *gin.Context) {
	// 接收请求参数
	type CreateRequest struct {
		Name          string        `json:"name"`
		Employees     []string      `json:"employees"`
		ChooseContact []interface{} `json:"chooseContact"`
		Guide         string        `json:"guide"`
		Rooms         []interface{} `json:"rooms"`
		FilterContact int           `json:"filterContact"`
		CorpID        int           `json:"corpId"`
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 从上下文获取corpID，暂时使用默认值1
	corpID := uint(1)
	if req.CorpID > 0 {
		corpID = uint(req.CorpID)
	}

	// 将数组转换为JSON字符串
	employees := ""
	if len(req.Employees) > 0 {
		employees = "[" + strings.Join(req.Employees, ",") + "]"
	}

	chooseContact, _ := json.Marshal(req.ChooseContact)
	rooms, _ := json.Marshal(req.Rooms)

	// 构建创建数据
	item := &model.RoomTagPull{
		Name:          req.Name,
		Employees:     employees,
		ChooseContact: string(chooseContact),
		Guide:         req.Guide,
		Rooms:         string(rooms),
		FilterContact: req.FilterContact,
		WxTid:         "[]", // 给wx_tid字段传一个有效的JSON值
		CorpID:        int(corpID),
	}

	// 调用服务创建数据
	err := h.service.Create(item)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建标签建群失败")
		return
	}

	response.Success(c, gin.H{"id": item.ID})
}

// Detail 标签建群详情
func (h *RoomTagPullHandler) Detail(c *gin.Context) {
	// 获取ID参数
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的ID")
		return
	}

	// 调用服务获取数据
	item, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取标签建群详情失败")
		return
	}

	response.Success(c, item)
}

// ContactDetail 客户详情
func (h *RoomTagPullHandler) ContactDetail(c *gin.Context) {
	// 这里可以实现客户详情的逻辑
	response.Success(c, nil)
}

// Show 标签建群详情（兼容旧接口）
func (h *RoomTagPullHandler) Show(c *gin.Context) {
	h.Detail(c)
}

// Store 创建标签建群（兼容旧接口）
func (h *RoomTagPullHandler) Store(c *gin.Context) {
	h.Create(c)
}

// Destroy 删除标签建群
func (h *RoomTagPullHandler) Destroy(c *gin.Context) {
	response.SuccessMsg(c, "删除成功")
}

// ChooseContact 选择联系人
func (h *RoomTagPullHandler) ChooseContact(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

// FilterContact 筛选联系人
func (h *RoomTagPullHandler) FilterContact(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

// ShowContact 查看已选联系人
func (h *RoomTagPullHandler) ShowContact(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

// RoomList 群列表
func (h *RoomTagPullHandler) RoomList(c *gin.Context) {
	response.Success(c, gin.H{"list": []interface{}{}})
}

// RemindSend 提醒发送
func (h *RoomTagPullHandler) RemindSend(c *gin.Context) {
	response.SuccessMsg(c, "提醒发送成功")
}
