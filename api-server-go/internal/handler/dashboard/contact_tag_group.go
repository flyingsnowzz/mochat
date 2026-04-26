// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含客户标签组管理的处理器：
// WorkContactTagGroupHandler - 处理客户标签组的增删改查操作
package dashboard

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

// WorkContactTagGroupHandler 客户标签组管理处理器
// 提供客户标签组的增删改查功能
// 主要职责：
// 1. 获取标签组列表
// 2. 获取标签组详情
// 3. 创建标签组
// 4. 更新标签组
// 5. 删除标签组
//
// 依赖服务：
// - WorkContactTagGroupService: 客户标签组服务

type WorkContactTagGroupHandler struct {
	groupSvc *service.WorkContactTagGroupService // 客户标签组服务
}

// NewWorkContactTagGroupHandler 创建客户标签组管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：客户标签组管理处理器实例
func NewWorkContactTagGroupHandler(db *gorm.DB) *WorkContactTagGroupHandler {
	return &WorkContactTagGroupHandler{
		groupSvc: service.NewWorkContactTagGroupService(db),
	}
}

// Index 获取标签组列表
// 获取企业的标签组列表
// 处理流程：
// 1. 获取企业 ID
// 2. 调用服务层获取标签组列表
// 3. 返回标签组列表
//
// 返回：包含标签组列表的响应
func (h *WorkContactTagGroupHandler) Index(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		response.Fail(c, response.ErrParams, "未获取到企业信息")
		return
	}

	// 调用服务层获取标签组列表
	groups, err := h.groupSvc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取标签组列表失败")
		return
	}

	// 返回标签组列表
	response.Success(c, gin.H{"list": groups})
}

// Detail 获取标签组详情
// 根据标签组 ID 获取其详细信息
// 处理流程：
// 1. 从路径参数中获取标签组 ID
// 2. 调用服务层获取标签组详情
// 3. 返回标签组详情信息
// 参数：
//
//	id - 标签组 ID（路径参数）
//
// 返回：标签组详情信息
func (h *WorkContactTagGroupHandler) Detail(c *gin.Context) {
	// 从路径参数中获取标签组 ID
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	// 调用服务层获取标签组详情
	group, err := h.groupSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "标签组不存在")
		return
	}

	// 返回标签组详情信息
	response.Success(c, group)
}

// Store 创建标签组
// 创建新的标签组
// 处理流程：
// 1. 获取企业 ID
// 2. 绑定请求参数到标签组对象
// 3. 设置标签组的企业 ID
// 4. 调用服务层创建标签组
// 5. 返回创建的标签组信息
// 请求体（JSON）：
//
//	包含标签组的相关信息，如名称等
//
// 返回：创建的标签组详情信息
func (h *WorkContactTagGroupHandler) Store(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		response.Fail(c, response.ErrParams, "未获取到企业信息")
		return
	}

	// 绑定请求参数到标签组对象
	var group model.WorkContactTagGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 设置标签组的企业 ID
	group.CorpID = corpID.(uint)
	
	// 调用服务层创建标签组
	if err := h.groupSvc.Create(&group); err != nil {
		response.Fail(c, response.ErrDB, "创建标签组失败")
		return
	}

	// 返回创建的标签组信息
	response.Success(c, group)
}

// Update 更新标签组
// 更新指定标签组的信息
// 处理流程：
// 1. 从路径参数中获取标签组 ID
// 2. 调用服务层获取标签组详情
// 3. 绑定请求参数到标签组对象
// 4. 调用服务层更新标签组
// 5. 返回更新后的标签组信息
// 参数：
//
//	id - 标签组 ID（路径参数）
//
// 请求体（JSON）：
//
//	包含标签组的相关信息，如名称等
//
// 返回：更新后的标签组详情信息
func (h *WorkContactTagGroupHandler) Update(c *gin.Context) {
	// 从路径参数中获取标签组 ID
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	// 调用服务层获取标签组详情
	group, err := h.groupSvc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "标签组不存在")
		return
	}

	// 绑定请求参数到标签组对象
	if err := c.ShouldBindJSON(group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}

	// 调用服务层更新标签组
	if err := h.groupSvc.Update(group); err != nil {
		response.Fail(c, response.ErrDB, "更新标签组失败")
		return
	}

	// 返回更新后的标签组信息
	response.Success(c, group)
}

// Destroy 删除标签组
// 删除指定的标签组
// 处理流程：
// 1. 从路径参数中获取标签组 ID
// 2. 调用服务层删除标签组
// 3. 返回删除结果
// 参数：
//
//	id - 标签组 ID（路径参数）
//
// 返回：删除成功的消息
func (h *WorkContactTagGroupHandler) Destroy(c *gin.Context) {
	// 从路径参数中获取标签组 ID
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	// 调用服务层删除标签组
	if err := h.groupSvc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除标签组失败")
		return
	}

	// 返回删除结果
	response.SuccessMsg(c, "删除成功")
}
