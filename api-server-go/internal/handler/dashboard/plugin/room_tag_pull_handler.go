package plugin

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

	// 调用服务获取数据
	result, err := h.service.List(name, page, perPage)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取标签建群列表失败")
		return
	}

	response.Success(c, result)
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
	corpID := 1
	if req.CorpID > 0 {
		corpID = req.CorpID
	}

	// 将数组转换为JSON字符串
	employees := ""
	if len(req.Employees) > 0 {
		employees = "[" + strings.Join(req.Employees, ",") + ""
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
		CorpID:        corpID,
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
