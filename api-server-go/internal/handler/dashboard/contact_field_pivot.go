// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含客户字段值管理的处理器：
// ContactFieldPivotHandler - 处理客户字段值的查询和更新操作
package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

// ContactFieldPivotHandler 客户字段值管理处理器
// 提供客户字段值的查询和更新功能
// 主要职责：
// 1. 获取客户的字段值列表
// 2. 批量更新客户的字段值
//
// 依赖服务：
// - ContactFieldPivotService: 客户字段值关联服务

type ContactFieldPivotHandler struct {
	pivotSvc *service.ContactFieldPivotService // 客户字段值关联服务
}

// NewContactFieldPivotHandler 创建客户字段值管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：客户字段值管理处理器实例
func NewContactFieldPivotHandler(db *gorm.DB) *ContactFieldPivotHandler {
	return &ContactFieldPivotHandler{
		pivotSvc: service.NewContactFieldPivotService(db),
	}
}

// Index 获取客户字段值列表
// 根据客户 ID 获取其所有字段的值列表
// 处理流程：
// 1. 获取 contactId 查询参数
// 2. 解析 contactId 为 uint 类型
// 3. 调用服务层获取字段值列表
// 4. 返回字段值列表
// 参数：
//
//	contactId - 客户 ID（查询参数）
//
// 返回：包含字段值列表的响应
func (h *ContactFieldPivotHandler) Index(c *gin.Context) {
	// 获取 contactId 查询参数
	contactIDStr := c.Query("contactId")
	if contactIDStr == "" {
		response.Fail(c, response.ErrParams, "缺少 contactId 参数")
		return
	}

	// 解析 contactId 为 uint 类型
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		response.Fail(c, response.ErrParams, "contactId 参数错误")
		return
	}

	// 调用服务层获取字段值列表
	pivots, err := h.pivotSvc.List(uint(contactID))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取字段值列表失败")
		return
	}

	// 返回字段值列表
	response.Success(c, gin.H{"list": pivots})
}

// Update 批量更新客户字段值
// 批量更新指定客户的字段值
// 处理流程：
// 1. 从路径参数中获取客户 ID
// 2. 绑定请求参数
// 3. 调用服务层批量更新字段值
// 4. 返回更新结果
// 参数：
//
//	id - 客户 ID（路径参数）
//
// 请求体（JSON）：
//
//	fields - 字段值列表，包含 fieldId 和 value
//
// 返回：更新成功的消息
func (h *ContactFieldPivotHandler) Update(c *gin.Context) {
	// 从路径参数中获取客户 ID
	contactID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 绑定请求参数
	var req struct {
		Fields []struct {
			FieldID uint   `json:"fieldId"` // 字段 ID
			Value   string `json:"value"`   // 字段值
		}
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 调用服务层批量更新字段值
	if err := h.pivotSvc.BatchUpdate(uint(contactID), req.Fields); err != nil {
		response.Fail(c, response.ErrDB, "更新字段值失败")
		return
	}

	// 返回更新结果
	response.SuccessMsg(c, "更新成功")
}
