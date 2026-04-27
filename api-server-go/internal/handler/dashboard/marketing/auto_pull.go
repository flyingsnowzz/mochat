package marketing

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type WorkRoomAutoPullHandler struct {
	service *plugin.WorkRoomAutoPullService
}

func NewWorkRoomAutoPullHandler(service *plugin.WorkRoomAutoPullService) *WorkRoomAutoPullHandler {
	return &WorkRoomAutoPullHandler{
		service: service,
	}
}

// Index 自动拉群列表
func (h *WorkRoomAutoPullHandler) Index(c *gin.Context) {
	// 获取查询参数
	qrcodeName := c.Query("qrcodeName")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))

	// 计算offset和limit
	offset := (page - 1) * perPage
	limit := perPage

	// 从上下文获取corpID，暂时使用默认值1
	corpID := uint(1)

	// 调用服务获取数据
	items, total, err := h.service.List(corpID, qrcodeName, offset, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取自动拉群列表失败")
		return
	}

	// 构建响应数据
	result := map[string]interface{}{
		"list":  items,
		"total": total,
		"page": map[string]int{
			"currentPage": page,
			"perPage":     perPage,
			"total":       int(total),
			"totalPage":   (int(total) + perPage - 1) / perPage,
		},
		"pageSize": perPage,
	}

	response.Success(c, result)
}

// Show 自动拉群详情
func (h *WorkRoomAutoPullHandler) Show(c *gin.Context) {
	// 获取ID参数
	idStr := c.Query("workRoomAutoPullId")
	// 兼容旧的参数名
	if idStr == "" {
		idStr = c.Query("id")
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的ID")
		return
	}

	// 调用服务获取数据
	item, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取自动拉群详情失败")
		return
	}

	// 转换为前端期望的格式
	type ShowResponse struct {
		QrcodeName   string        `json:"qrcodeName"`
		IsVerified   int           `json:"isVerified"`
		LeadingWords string        `json:"leadingWords"`
		Employees    []interface{} `json:"employees"`
		Tags         []interface{} `json:"tags"`
		Rooms        []interface{} `json:"rooms"`
		SelectedTags []string      `json:"selectedTags"`
	}

	// 解析employees字段（逗号分隔的字符串转换为数组）
	employees := []interface{}{}
	if item.Employees != "" {
		// 这里需要根据实际存储格式解析，暂时返回空数组
	}

	// 解析tags字段（JSON字符串转换为数组）
	tags := []interface{}{}
	if item.Tags != "" {
		// 这里需要根据实际存储格式解析，暂时返回空数组
	}

	// 解析rooms字段（JSON字符串转换为数组）
	rooms := []interface{}{}
	if item.Rooms != "" {
		// 这里需要根据实际存储格式解析，暂时返回空数组
	}

	// 构建响应数据
	responseData := ShowResponse{
		QrcodeName:   item.QrcodeName,
		IsVerified:   item.IsVerified,
		LeadingWords: item.LeadingWords,
		Employees:    employees,
		Tags:         tags,
		Rooms:        rooms,
		SelectedTags: []string{}, // 暂时返回空数组
	}

	response.Success(c, responseData)
}

// Store 创建自动拉群
func (h *WorkRoomAutoPullHandler) Store(c *gin.Context) {
	// 接收请求参数
	type StoreRequest struct {
		QrcodeName   string `json:"qrcodeName"`
		IsVerified   int    `json:"isVerified"`
		LeadingWords string `json:"leadingWords"`
		Employees    string `json:"employees"`
		Tags         string `json:"tags"`
		Rooms        string `json:"rooms"`
		CorpID       uint   `json:"corpId"`
	}

	var req StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 从上下文获取corpID，暂时使用默认值1
	corpID := uint(1)
	if req.CorpID > 0 {
		corpID = req.CorpID
	}

	// 创建自动拉群
	item := &model.WorkRoomAutoPull{
		CorpID:       corpID,
		QrcodeName:   req.QrcodeName,
		IsVerified:   req.IsVerified,
		LeadingWords: req.LeadingWords,
		Employees:    req.Employees,
		Tags:         req.Tags,
		Rooms:        req.Rooms,
	}

	// 调用服务创建数据
	err := h.service.Create(item)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建自动拉群失败")
		return
	}

	response.Success(c, gin.H{"id": item.ID})
}

// Update 更新自动拉群
func (h *WorkRoomAutoPullHandler) Update(c *gin.Context) {
	// 接收请求参数
	type UpdateRequest struct {
		ID           uint   `json:"id"`
		QrcodeName   string `json:"qrcodeName"`
		IsVerified   int    `json:"isVerified"`
		LeadingWords string `json:"leadingWords"`
		Employees    string `json:"employees"`
		Tags         string `json:"tags"`
		Rooms        string `json:"rooms"`
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 构建更新数据
	updates := map[string]interface{}{
		"qrcode_name":   req.QrcodeName,
		"is_verified":   req.IsVerified,
		"leading_words": req.LeadingWords,
		"employees":     req.Employees,
		"tags":          req.Tags,
		"rooms":         req.Rooms,
		"updated_at":    time.Now(),
	}

	// 调用服务更新数据
	err := h.service.Update(req.ID, updates)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新自动拉群失败")
		return
	}

	response.Success(c, nil)
}

// Destroy 删除自动拉群
func (h *WorkRoomAutoPullHandler) Destroy(c *gin.Context) {
	// 接收请求参数
	type DestroyRequest struct {
		ID uint `json:"id"`
	}

	var req DestroyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 尝试从查询参数获取ID
		idStr := c.Query("id")
		if idStr == "" {
			response.Fail(c, http.StatusBadRequest, "参数错误")
			return
		}
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "无效的ID")
			return
		}
		req.ID = uint(id)
	}

	// 调用服务删除数据
	err := h.service.Delete(req.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除自动拉群失败")
		return
	}

	response.Success(c, nil)
}
