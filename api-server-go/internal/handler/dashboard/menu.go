// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含菜单管理的处理器，提供菜单的 CRUD 操作和相关功能
package dashboard

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

// MenuHandler 菜单管理处理器
// 提供菜单的 CRUD 操作和相关功能
// 主要职责：
// 1. 获取菜单列表
// 2. 获取菜单选择列表
// 3. 获取菜单详情
// 4. 获取图标列表
// 5. 创建菜单
// 6. 更新菜单
// 7. 删除菜单
// 8. 更新菜单状态
//
// 依赖服务：
// - RbacMenuService: 菜单服务

type MenuHandler struct {
	svc *service.RbacMenuService // 菜单服务
}

// menuIndexItem 菜单列表项
// 用于构建菜单列表的树形结构
// 字段说明：
// - MenuID: 菜单 ID
// - MenuPath: 菜单路径
// - Name: 菜单名称
// - Level: 菜单级别
// - LevelName: 菜单级别名称
// - ParentID: 父菜单 ID
// - Icon: 菜单图标
// - Status: 菜单状态
// - OperateName: 操作人名称
// - UpdatedAt: 更新时间
// - Children: 子菜单列表

type menuIndexItem struct {
	MenuID      uint            `json:"menuId"`      // 菜单 ID
	MenuPath    string          `json:"menuPath"`    // 菜单路径
	Name        string          `json:"name"`        // 菜单名称
	Level       int             `json:"level"`       // 菜单级别
	LevelName   string          `json:"levelName"`   // 菜单级别名称
	ParentID    uint            `json:"parentId"`    // 父菜单 ID
	Icon        string          `json:"icon"`        // 菜单图标
	Status      string          `json:"status"`      // 菜单状态
	OperateName string          `json:"operateName"` // 操作人名称
	UpdatedAt   string          `json:"updatedAt"`   // 更新时间
	Children    []menuIndexItem `json:"children"`    // 子菜单列表
}

// NewMenuHandler 创建菜单管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：菜单管理处理器实例
func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{svc: service.NewRbacMenuService(db)}
}

// Index 获取菜单列表
// 获取菜单列表，支持按名称搜索，并返回树形结构
// 处理流程：
// 1. 获取搜索关键词
// 2. 调用服务层获取菜单列表
// 3. 构建菜单树形结构
// 4. 返回分页结果
// 参数：
//
//	name - 菜单名称（用于搜索）
//
// 返回：包含菜单列表、总数、分页信息的响应
func (h *MenuHandler) Index(c *gin.Context) {
	// 获取搜索关键词
	name := strings.TrimSpace(c.Query("name"))
	
	// 调用服务层获取菜单列表
	menus, err := h.svc.List(name)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取菜单列表失败")
		return
	}
	
	// 构建菜单树形结构
	items := buildMenuIndexTree(menus)
	
	// 返回分页结果
	response.PageResult(c, items, int64(len(items)), 1, len(items))
}

// Select 获取菜单选择列表
// 获取菜单选择列表，返回树形结构，用于选择菜单
// 处理流程：
// 1. 调用服务层获取菜单选择列表
// 2. 构建菜单选择树形结构
// 3. 返回菜单选择列表
//
// 返回：包含菜单选择列表的响应
func (h *MenuHandler) Select(c *gin.Context) {
	// 调用服务层获取菜单选择列表
	menus, err := h.svc.Select()
	if err != nil {
		response.Fail(c, response.ErrDB, "获取菜单选择列表失败")
		return
	}
	
	// 构建菜单选择树形结构并返回
	response.Success(c, buildMenuSelectTree(menus))
}

// Show 获取菜单详情
// 根据菜单 ID 获取菜单详情
// 处理流程：
// 1. 获取菜单 ID
// 2. 调用服务层获取菜单详情
// 3. 构建菜单详情数据
// 4. 返回菜单详情
//
// 返回：包含菜单详情的响应
func (h *MenuHandler) Show(c *gin.Context) {
	// 获取菜单 ID
	id := getMenuID(c)
	
	// 调用服务层获取菜单详情
	menu, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "菜单不存在")
		return
	}
	
	// 构建菜单详情数据并返回
	response.Success(c, buildMenuShowData(menu))
}

// IconIndex 获取图标列表
// 获取菜单图标列表（暂未实现）
// 处理流程：
// 1. 直接返回空数组
//
// 返回：空数组
func (h *MenuHandler) IconIndex(c *gin.Context) {
	response.Success(c, []string{})
}

// Store 创建菜单
// 创建新的菜单
// 处理流程：
// 1. 绑定请求参数
// 2. 调用服务层创建菜单
// 3. 返回创建的菜单
// 请求体（JSON）：
//
//	包含菜单的相关信息，如名称、图标、链接等
//
// 返回：创建的菜单详情
func (h *MenuHandler) Store(c *gin.Context) {
	// 绑定请求参数
	var menu model.RbacMenu
	if err := c.ShouldBindJSON(&menu); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层创建菜单
	if err := h.svc.Create(&menu); err != nil {
		response.Fail(c, response.ErrDB, "创建菜单失败")
		return
	}
	
	// 返回创建的菜单
	response.Success(c, menu)
}

// Update 更新菜单
// 更新现有菜单的信息
// 处理流程：
// 1. 获取菜单 ID
// 2. 调用服务层获取菜单详情
// 3. 绑定请求参数
// 4. 更新菜单字段（名称、图标、链接 URL、链接类型、是否页面菜单、数据权限）
// 5. 调用服务层更新菜单
// 6. 构建菜单详情数据
// 7. 返回更新后的菜单
// 请求体（JSON）：
//
//	menuId - 菜单 ID
//	name - 菜单名称
//	icon - 菜单图标
//	linkUrl - 链接 URL
//	linkType - 链接类型
//	isPageMenu - 是否页面菜单
//	dataPermission - 数据权限
//
// 返回：更新后的菜单详情
func (h *MenuHandler) Update(c *gin.Context) {
	// 获取菜单 ID
	id := getMenuID(c)
	
	// 调用服务层获取菜单详情
	menu, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "菜单不存在")
		return
	}
	
	// 绑定请求参数
	var req struct {
		MenuID         uint   `json:"menuId"`         // 菜单 ID
		Name           string `json:"name"`           // 菜单名称
		Icon           string `json:"icon"`           // 菜单图标
		LinkURL        string `json:"linkUrl"`        // 链接 URL
		LinkType       int    `json:"linkType"`       // 链接类型
		IsPageMenu     int    `json:"isPageMenu"`     // 是否页面菜单
		DataPermission int    `json:"dataPermission"` // 数据权限
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 更新菜单字段
	menu.Name = req.Name
	menu.Icon = req.Icon
	menu.LinkURL = req.LinkURL
	if req.LinkType > 0 {
		menu.LinkType = req.LinkType
	}
	if req.IsPageMenu > 0 {
		menu.IsPageMenu = req.IsPageMenu
	}
	if req.DataPermission > 0 {
		menu.DataPermission = req.DataPermission
	}
	
	// 调用服务层更新菜单
	if err := h.svc.Update(menu); err != nil {
		response.Fail(c, response.ErrDB, "更新菜单失败")
		return
	}
	
	// 构建菜单详情数据并返回
	response.Success(c, buildMenuShowData(menu))
}

// Destroy 删除菜单
// 删除指定的菜单
// 处理流程：
// 1. 获取菜单 ID
// 2. 调用服务层删除菜单
// 3. 返回删除结果
//
// 返回：删除成功的消息
func (h *MenuHandler) Destroy(c *gin.Context) {
	// 获取菜单 ID
	id := getMenuID(c)
	
	// 调用服务层删除菜单
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除菜单失败")
		return
	}
	
	// 返回删除结果
	response.SuccessMsg(c, "删除成功")
}

// StatusUpdate 更新菜单状态
// 更新指定菜单的状态
// 处理流程：
// 1. 绑定请求参数
// 2. 调用服务层获取菜单详情
// 3. 更新菜单状态
// 4. 调用服务层更新菜单
// 5. 返回状态更新结果
// 请求体（JSON）：
//
//	menuId - 菜单 ID（必填）
//	status - 菜单状态（必填）
//
// 返回：状态更新成功的消息
func (h *MenuHandler) StatusUpdate(c *gin.Context) {
	// 绑定请求参数
	var req struct {
		MenuID uint `json:"menuId" binding:"required"` // 菜单 ID
		Status int  `json:"status" binding:"required"`  // 菜单状态
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层获取菜单详情
	menu, err := h.svc.GetByID(req.MenuID)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "菜单不存在")
		return
	}
	
	// 更新菜单状态
	menu.Status = req.Status
	
	// 调用服务层更新菜单
	if err := h.svc.Update(menu); err != nil {
		response.Fail(c, response.ErrDB, "更新状态失败")
		return
	}
	
	// 返回状态更新结果
	response.SuccessMsg(c, "状态更新成功")
}

// buildMenuIndexTree 构建菜单列表树形结构
// 将菜单列表构建成树形结构，用于菜单列表展示
// 处理流程：
// 1. 构建父菜单 ID 到子菜单列表的映射
// 2. 递归构建菜单树形结构，从根菜单开始
// 参数：
//
//	menus - 菜单列表
//
// 返回：菜单列表树形结构
func buildMenuIndexTree(menus []model.RbacMenu) []menuIndexItem {
	// 构建父菜单 ID 到子菜单列表的映射
	childrenMap := make(map[uint][]model.RbacMenu)
	for _, menu := range menus {
		childrenMap[menu.ParentID] = append(childrenMap[menu.ParentID], menu)
	}
	
	// 递归构建菜单树形结构，从根菜单开始
	return buildMenuIndexChildren(childrenMap, 0, "")
}

// buildMenuIndexChildren 构建菜单列表子项树形结构
// 递归构建菜单列表的子项树形结构
// 处理流程：
// 1. 获取指定父菜单的子菜单列表
// 2. 遍历子菜单，构建菜单列表项
// 3. 为每个菜单列表项递归构建子菜单树形结构
// 4. 返回菜单列表项数组
// 参数：
//
//	childrenMap - 父菜单 ID 到子菜单列表的映射
//	parentID - 父菜单 ID
//	prefix - 路径前缀
//
// 返回：菜单列表项数组
func buildMenuIndexChildren(childrenMap map[uint][]model.RbacMenu, parentID uint, prefix string) []menuIndexItem {
	// 获取指定父菜单的子菜单列表
	children := childrenMap[parentID]
	items := make([]menuIndexItem, 0, len(children))
	
	// 遍历子菜单，构建菜单列表项
	for idx, menu := range children {
		// 构建菜单路径
		path := strconv.Itoa(idx + 1)
		if prefix != "" {
			path = prefix + "-" + path
		}
		
		// 格式化更新时间
		updatedAt := ""
		if !menu.UpdatedAt.IsZero() {
			updatedAt = menu.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		
		// 构建菜单列表项
		item := menuIndexItem{
			MenuID:      menu.ID,      // 菜单 ID
			MenuPath:    path,         // 菜单路径
			Name:        menu.Name,    // 菜单名称
			Level:       menu.Level,   // 菜单级别
			LevelName:   strconv.Itoa(menu.Level), // 菜单级别名称
			ParentID:    menu.ParentID, // 父菜单 ID
			Icon:        menu.Icon,    // 菜单图标
			Status:      strconv.Itoa(menu.Status), // 菜单状态
			OperateName: menu.OperateName, // 操作人名称
			UpdatedAt:   updatedAt,    // 更新时间
		}
		
		// 递归构建子菜单树形结构
		item.Children = buildMenuIndexChildren(childrenMap, menu.ID, path)
		items = append(items, item)
	}
	
	// 返回菜单列表项数组
	return items
}

// buildMenuSelectTree 构建菜单选择树形结构
// 将菜单列表构建成树形结构，用于菜单选择
// 处理流程：
// 1. 构建父菜单 ID 到子菜单列表的映射
// 2. 定义递归构建函数
// 3. 递归构建菜单树形结构，从根菜单开始
// 4. 返回菜单选择树形结构
// 参数：
//
//	menus - 菜单列表
//
// 返回：菜单选择树形结构
func buildMenuSelectTree(menus []model.RbacMenu) []gin.H {
	// 构建父菜单 ID 到子菜单列表的映射
	childrenMap := make(map[uint][]model.RbacMenu)
	for _, menu := range menus {
		childrenMap[menu.ParentID] = append(childrenMap[menu.ParentID], menu)
	}
	
	// 定义递归构建函数
	var build func(parentID uint) []gin.H
	build = func(parentID uint) []gin.H {
		// 获取指定父菜单的子菜单列表
		children := childrenMap[parentID]
		items := make([]gin.H, 0, len(children))
		
		// 遍历子菜单，构建菜单选择项
		for _, menu := range children {
			items = append(items, gin.H{
				"menuId":         menu.ID,         // 菜单 ID
				"name":           menu.Name,       // 菜单名称
				"level":          menu.Level,      // 菜单级别
				"dataPermission": menu.DataPermission, // 数据权限
				"parentId":       menu.ParentID,   // 父菜单 ID
				"children":       build(menu.ID),  // 子菜单列表
			})
		}
		return items
	}
	
	// 递归构建菜单树形结构，从根菜单开始
	return build(0)
}

// buildMenuShowData 构建菜单详情数据
// 构建菜单详情数据，用于菜单编辑和展示
// 处理流程：
// 1. 构建包含菜单详情的 gin.H 数据
// 2. 从菜单路径中提取各级父菜单 ID
// 3. 添加菜单的其他属性
// 参数：
//
//	menu - 菜单实例
//
// 返回：菜单详情数据
func buildMenuShowData(menu *model.RbacMenu) gin.H {
	return gin.H{
		"menuId":         menu.ID,         // 菜单 ID
		"firstMenuId":    firstParentFromPath(menu.Path, 1),    // 一级父菜单 ID
		"secondMenuId":   firstParentFromPath(menu.Path, 2),   // 二级父菜单 ID
		"thirdMenuId":    firstParentFromPath(menu.Path, 3),    // 三级父菜单 ID
		"fourthMenuId":   firstParentFromPath(menu.Path, 4),   // 四级父菜单 ID
		"level":          menu.Level,      // 菜单级别
		"name":           menu.Name,       // 菜单名称
		"icon":           menu.Icon,       // 菜单图标
		"isPageMenu":     menu.IsPageMenu,     // 是否页面菜单
		"linkUrl":        menu.LinkURL,    // 链接 URL
		"linkType":       menu.LinkType,   // 链接类型
		"dataPermission": menu.DataPermission, // 数据权限
		"status":         menu.Status,     // 菜单状态
	}
}

// firstParentFromPath 从路径中获取指定级别的父菜单 ID
// 从菜单路径中提取指定级别的父菜单 ID
// 处理流程：
// 1. 如果路径为空，返回 0
// 2. 按 "#-#" 分割路径
// 3. 如果分割后的部分数量小于等于指定索引，返回 0
// 4. 去除指定索引部分的 "#" 字符
// 5. 将字符串转换为 uint 类型并返回
// 参数：
//
//	path - 菜单路径
//	index - 索引位置
//
// 返回：指定级别的父菜单 ID 或 0
func firstParentFromPath(path string, index int) uint {
	// 如果路径为空，返回 0
	if path == "" {
		return 0
	}
	
	// 按 "#-#" 分割路径
	parts := strings.Split(path, "#-#")
	
	// 如果分割后的部分数量小于等于指定索引，返回 0
	if len(parts) <= index {
		return 0
	}
	
	// 去除指定索引部分的 "#" 字符
	value := strings.Trim(parts[index], "#")
	
	// 将字符串转换为 uint 类型并返回
	id, _ := strconv.ParseUint(value, 10, 32)
	return uint(id)
}

// getMenuID 获取菜单 ID
// 从 Gin 上下文中获取菜单 ID，支持从路径参数、查询参数和请求体中获取
// 处理流程：
// 1. 尝试从路径参数中获取 ID
// 2. 如果路径参数中没有，尝试从查询参数中获取
// 3. 如果查询参数中没有，尝试从请求体中获取
// 4. 如果都没有，返回 0
// 参数：
//
//	c - Gin 上下文
//
// 返回：菜单 ID 或 0
func getMenuID(c *gin.Context) uint64 {
	// 尝试从路径参数中获取 ID
	if id, err := strconv.ParseUint(c.Param("id"), 10, 32); err == nil && id > 0 {
		return id
	}
	
	// 尝试从查询参数中获取 ID
	if id, err := strconv.ParseUint(c.Query("menuId"), 10, 32); err == nil && id > 0 {
		return id
	}
	
	// 尝试从请求体中获取 ID
	var req struct {
		MenuID uint `json:"menuId"` // 菜单 ID
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.MenuID > 0 {
		return uint64(req.MenuID)
	}
	
	// 如果都没有，返回 0
	return 0
}
