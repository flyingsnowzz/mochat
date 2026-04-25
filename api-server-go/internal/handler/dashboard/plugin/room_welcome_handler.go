package plugin

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service/plugin"
)

type RoomWelcomeHandler struct {
	service *plugin.RoomWelcomeService
}

func NewRoomWelcomeHandler(service *plugin.RoomWelcomeService) *RoomWelcomeHandler {
	return &RoomWelcomeHandler{service: service}
}

func (h *RoomWelcomeHandler) Index(c *gin.Context) {
	// 从上下文获取企业ID
	corpID := uint(1) // 这里应该从认证中间件获取实际的企业ID
	
	// 调用服务获取数据
	items, err := h.service.List(corpID)
	if err != nil {
		response.Fail(c, 500, "获取入群欢迎语列表失败")
		return
	}
	
	// 转换为响应格式
	var list []interface{}
	for _, item := range items {
		list = append(list, item)
	}
	
	response.Success(c, gin.H{"list": list, "page": gin.H{"total": len(list)}})
}

func (h *RoomWelcomeHandler) Show(c *gin.Context) {
	// 获取ID参数（使用查询参数而不是路径参数）
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}

	// 调用服务获取数据
	item, err := h.service.GetByID(uint(id))
	if err != nil {
		response.Fail(c, 500, "获取入群欢迎语详情失败")
		return
	}

	response.Success(c, item)
}

func (h *RoomWelcomeHandler) Store(c *gin.Context) {
	// 绑定请求数据
	var req struct {
		MsgText     string      `json:"msg_text"`
		MsgComplex  interface{} `json:"msg_complex"`
		ComplexType int         `json:"complex_type"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "无效的请求数据")
		return
	}
	
	// 从上下文获取企业ID
	corpID := uint(1) // 这里应该从认证中间件获取实际的企业ID
	
	// 处理 MsgComplex，将 interface{} 转换为 string
	var msgComplexStr string
	if req.MsgComplex != nil {
		msgComplexBytes, err := json.Marshal(req.MsgComplex)
		if err != nil {
			response.Fail(c, 400, "无效的消息复杂内容")
			return
		}
		msgComplexStr = string(msgComplexBytes)
	}
	
	// 创建模型
	item := &model.RoomWelcomeTemplate{
		CorpID:      corpID,
		MsgText:     req.MsgText,
		MsgComplex:  msgComplexStr,
		ComplexType: req.ComplexType,
	}
	
	// 调用服务保存数据
	if err := h.service.Create(item); err != nil {
		response.Fail(c, 500, "创建入群欢迎语失败")
		return
	}
	
	response.SuccessMsg(c, "创建成功")
}

func (h *RoomWelcomeHandler) Update(c *gin.Context) {
	// 绑定请求数据
	var req struct {
		ID          string      `json:"id"`
		MsgText     string      `json:"msg_text"`
		MsgComplex  interface{} `json:"msg_complex"`
		ComplexType int         `json:"complex_type"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "无效的请求数据")
		return
	}
	
	// 获取ID参数
	id, err := strconv.ParseUint(req.ID, 10, 32)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	
	// 处理 MsgComplex，将 interface{} 转换为 string
	var msgComplexStr string
	if req.MsgComplex != nil {
		msgComplexBytes, err := json.Marshal(req.MsgComplex)
		if err != nil {
			response.Fail(c, 400, "无效的消息复杂内容")
			return
		}
		msgComplexStr = string(msgComplexBytes)
	}
	
	// 准备更新数据
	updates := map[string]interface{}{
		"msg_text":      req.MsgText,
		"msg_complex":   msgComplexStr,
		"complex_type":  req.ComplexType,
	}
	
	// 调用服务更新数据
	if err := h.service.Update(uint(id), updates); err != nil {
		response.Fail(c, 500, "更新入群欢迎语失败")
		return
	}
	
	response.SuccessMsg(c, "更新成功")
}

func (h *RoomWelcomeHandler) Destroy(c *gin.Context) {
	// 绑定请求数据
	var req struct {
		ID string `json:"id"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "无效的请求数据")
		return
	}
	
	// 获取ID参数
	id, err := strconv.ParseUint(req.ID, 10, 32)
	if err != nil {
		response.Fail(c, 400, "无效的ID")
		return
	}
	
	// 调用服务删除数据
	if err := h.service.Delete(uint(id)); err != nil {
		response.Fail(c, 500, "删除入群欢迎语失败")
		return
	}
	
	response.SuccessMsg(c, "删除成功")
}
