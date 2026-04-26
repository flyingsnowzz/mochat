// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该文件包含素材和素材分组的处理器：
// 1. MediumHandler - 处理素材的 CRUD 操作
// 2. MediumGroupHandler - 处理素材分组的 CRUD 操作
package dashboard

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/service"
)

// MediumHandler 素材管理处理器
// 提供素材的 CRUD 操作和分组管理功能
// 主要职责：
// 1. 获取素材列表
// 2. 获取素材详情
// 3. 创建素材
// 4. 更新素材
// 5. 删除素材
// 6. 移动素材到指定分组
//
// 依赖服务：
// - MediumService: 素材服务
// - MediumGroupService: 素材分组服务
// - gorm.DB: 数据库连接

type MediumHandler struct {
	db       *gorm.DB                     // 数据库连接
	svc      *service.MediumService        // 素材服务
	groupSvc *service.MediumGroupService   // 素材分组服务
}

// NewMediumHandler 创建素材管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：素材管理处理器实例
func NewMediumHandler(db *gorm.DB) *MediumHandler {
	return &MediumHandler{
		db:       db,
		svc:      service.NewMediumService(db),
		groupSvc: service.NewMediumGroupService(db),
	}
}

// Index 获取素材列表
// 获取素材列表，支持分页、类型筛选、分组筛选和搜索
// 处理流程：
// 1. 获取企业 ID
// 2. 获取分页参数和筛选条件
// 3. 调用服务层获取素材列表
// 4. 加载素材分组名称
// 5. 构建返回数据
// 6. 返回分页结果
// 参数：
//
//	page - 页码，默认为 1
//	perPage/pageSize - 每页数量，默认为 10
//	type - 素材类型，默认为 0（全部）
//	mediumGroupId/groupId - 素材分组 ID，默认为 0（未分组）
//	searchStr - 搜索关键词
//
// 返回：包含素材列表、总数、分页信息的响应
func (h *MediumHandler) Index(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	
	// 获取分页参数和筛选条件
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))
	mediumType, _ := strconv.Atoi(c.DefaultQuery("type", "0"))
	groupID, _ := strconv.ParseUint(c.DefaultQuery("mediumGroupId", c.DefaultQuery("groupId", "0")), 10, 32)
	searchStr := c.Query("searchStr")

	// 调用服务层获取素材列表
	media, total, err := h.svc.ListWithSearch(corpID.(uint), page, pageSize, mediumType, uint(groupID), searchStr)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取素材列表失败")
		return
	}
	
	// 加载素材分组名称
	groupNames := h.loadMediumGroupNames(corpID.(uint))
	
	// 构建返回数据
	items := make([]gin.H, 0, len(media))
	for _, medium := range media {
		items = append(items, buildMediumResponseItem(medium, groupNames[medium.MediumGroupID]))
	}
	
	// 返回分页结果
	response.PageResult(c, items, total, page, pageSize)
}

// Show 获取素材详情
// 根据素材 ID 获取素材详情
// 处理流程：
// 1. 获取素材 ID
// 2. 调用服务层获取素材详情
// 3. 加载素材分组名称
// 4. 构建返回数据
// 5. 返回素材详情
// 参数：
//
//	id - 素材 ID（路径参数或查询参数）
//
// 返回：包含素材详情的响应
func (h *MediumHandler) Show(c *gin.Context) {
	// 获取素材 ID
	id := getUintID(c, "id")
	
	// 调用服务层获取素材详情
	medium, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "素材不存在")
		return
	}

	// 加载素材分组名称
	corpID, _ := c.Get("corpId")
	groupNames := h.loadMediumGroupNames(asUint(corpID))
	
	// 构建返回数据并返回
	response.Success(c, buildMediumResponseItem(*medium, groupNames[medium.MediumGroupID]))
}

// Store 创建素材
// 创建新的素材
// 处理流程：
// 1. 绑定请求参数
// 2. 获取企业 ID 和用户 ID
// 3. 验证企业 ID 是否有效
// 4. 序列化素材内容
// 5. 设置默认的 IsSync 值
// 6. 创建素材实例并设置相关字段
// 7. 调用服务层创建素材
// 8. 加载素材分组名称
// 9. 构建返回数据并返回
// 请求体（JSON）：
//
//	type - 素材类型（必填）
//	isSync - 是否同步，默认为 1
//	content - 素材内容（必填）
//	mediumGroupId - 素材分组 ID
//
// 返回：创建的素材详情
func (h *MediumHandler) Store(c *gin.Context) {
	// 绑定请求参数
	var req struct {
		Type          int                    `json:"type" binding:"required"`          // 素材类型
		IsSync        int                    `json:"isSync"`                          // 是否同步
		Content       map[string]interface{} `json:"content" binding:"required"`       // 素材内容
		MediumGroupID uint                   `json:"mediumGroupId"`                   // 素材分组 ID
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 获取企业 ID 和用户 ID
	corpID, _ := c.Get("corpId")
	userID, _ := c.Get("userId")
	if corpID == nil || corpID.(uint) == 0 {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}

	// 序列化素材内容
	content, err := service.MarshalMediumContent(req.Content)
	if err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 设置默认的 IsSync 值
	if req.IsSync == 0 {
		req.IsSync = 1
	}

	// 创建素材实例并设置相关字段
	medium := model.Medium{
		Type:          req.Type,          // 素材类型
		IsSync:        req.IsSync,        // 是否同步
		Content:       content,           // 素材内容
		CorpID:        corpID.(uint),     // 企业 ID
		MediumGroupID: req.MediumGroupID, // 素材分组 ID
		UserID:        asUint(userID),    // 用户 ID
		UserName:      h.loadUserName(asUint(userID)), // 用户名
	}
	
	// 调用服务层创建素材
	if err := h.svc.Create(&medium); err != nil {
		response.Fail(c, response.ErrDB, "创建素材失败")
		return
	}

	// 加载素材分组名称并返回
	groupNames := h.loadMediumGroupNames(corpID.(uint))
	response.Success(c, buildMediumResponseItem(medium, groupNames[medium.MediumGroupID]))
}

// Update 更新素材
// 更新现有素材的信息
// 处理流程：
// 1. 获取素材 ID
// 2. 调用服务层获取素材详情
// 3. 绑定请求参数
// 4. 更新素材字段（类型、是否同步、内容、分组）
// 5. 调用服务层更新素材
// 6. 加载素材分组名称
// 7. 构建返回数据并返回
// 请求体（JSON）：
//
//	id - 素材 ID
//	type - 素材类型
//	isSync - 是否同步
//	content - 素材内容
//	mediumGroupId - 素材分组 ID
//
// 返回：更新后的素材详情
func (h *MediumHandler) Update(c *gin.Context) {
	// 获取素材 ID
	id := getUintID(c, "id")
	
	// 调用服务层获取素材详情
	medium, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "素材不存在")
		return
	}

	// 绑定请求参数
	var req struct {
		ID            uint                   `json:"id"`            // 素材 ID
		Type          int                    `json:"type"`          // 素材类型
		IsSync        int                    `json:"isSync"`        // 是否同步
		Content       map[string]interface{} `json:"content"`       // 素材内容
		MediumGroupID uint                   `json:"mediumGroupId"` // 素材分组 ID
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 更新素材字段
	if req.Type > 0 {
		medium.Type = req.Type
	}
	if req.IsSync > 0 {
		medium.IsSync = req.IsSync
	}
	if req.Content != nil {
		content, err := service.MarshalMediumContent(req.Content)
		if err != nil {
			response.Fail(c, response.ErrParams, "参数错误")
			return
		}
		medium.Content = content
	}
	medium.MediumGroupID = req.MediumGroupID
	
	// 调用服务层更新素材
	if err := h.svc.Update(medium); err != nil {
		response.Fail(c, response.ErrDB, "更新素材失败")
		return
	}

	// 加载素材分组名称并返回
	corpID, _ := c.Get("corpId")
	groupNames := h.loadMediumGroupNames(asUint(corpID))
	response.Success(c, buildMediumResponseItem(*medium, groupNames[medium.MediumGroupID]))
}

// Destroy 删除素材
// 删除指定的素材
// 处理流程：
// 1. 获取素材 ID
// 2. 调用服务层删除素材
// 3. 返回删除结果
// 参数：
//
//	id - 素材 ID（路径参数或查询参数）
//
// 返回：删除成功的消息
func (h *MediumHandler) Destroy(c *gin.Context) {
	// 获取素材 ID
	id := getUintID(c, "id")
	
	// 调用服务层删除素材
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除素材失败")
		return
	}
	
	// 返回删除结果
	response.SuccessMsg(c, "删除成功")
}

// GroupUpdate 移动素材分组
// 将素材移动到指定的分组
// 处理流程：
// 1. 绑定请求参数
// 2. 调用服务层获取素材详情
// 3. 更新素材的分组 ID
// 4. 调用服务层更新素材
// 5. 返回移动结果
// 请求体（JSON）：
//
//	id - 素材 ID（必填）
//	mediumGroupId - 目标分组 ID
//
// 返回：移动成功的消息
func (h *MediumHandler) GroupUpdate(c *gin.Context) {
	// 绑定请求参数
	var req struct {
		ID      uint `json:"id" binding:"required"`      // 素材 ID
		GroupID uint `json:"mediumGroupId"`              // 目标分组 ID
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层获取素材详情
	medium, err := h.svc.GetByID(req.ID)
	if err != nil {
		response.Fail(c, response.ErrNotFound, "素材不存在")
		return
	}
	
	// 更新素材的分组 ID
	medium.MediumGroupID = req.GroupID
	
	// 调用服务层更新素材
	if err := h.svc.Update(medium); err != nil {
		response.Fail(c, response.ErrDB, "移动素材失败")
		return
	}
	
	// 返回移动结果
	response.SuccessMsg(c, "移动成功")
}

// loadMediumGroupNames 加载素材分组名称
// 加载指定企业的素材分组名称，返回分组 ID 到名称的映射
// 处理流程：
// 1. 初始化结果映射，默认包含 "未分组" 项
// 2. 如果企业 ID 为 0，直接返回默认结果
// 3. 调用服务层获取分组列表
// 4. 构建分组 ID 到名称的映射
// 5. 返回映射结果
// 参数：
//
//	corpID - 企业 ID
//
// 返回：分组 ID 到名称的映射
func (h *MediumHandler) loadMediumGroupNames(corpID uint) map[uint]string {
	// 初始化结果映射，默认包含 "未分组" 项
	result := map[uint]string{0: "未分组"}
	if corpID == 0 {
		return result
	}
	
	// 调用服务层获取分组列表
	groups, err := h.groupSvc.List(corpID)
	if err != nil {
		return result
	}
	
	// 构建分组 ID 到名称的映射
	for _, group := range groups {
		result[group.ID] = group.Name
	}
	
	// 返回映射结果
	return result
}

// loadUserName 加载用户名
// 根据用户 ID 加载用户名
// 处理流程：
// 1. 如果用户 ID 为 0，返回空字符串
// 2. 创建用户服务实例
// 3. 调用服务层获取用户详情
// 4. 返回用户名称
// 参数：
//
//	userID - 用户 ID
//
// 返回：用户名称或空字符串
func (h *MediumHandler) loadUserName(userID uint) string {
	// 如果用户 ID 为 0，返回空字符串
	if userID == 0 {
		return ""
	}
	
	// 创建用户服务实例
	userSvc := service.NewUserService(h.db)
	
	// 调用服务层获取用户详情
	user, err := userSvc.GetByID(userID)
	if err != nil {
		return ""
	}
	
	// 返回用户名称
	return user.Name
}

// buildMediumResponseItem 构建素材响应项
// 构建素材的响应数据，包括解析内容、添加完整路径、生成标题等
// 处理流程：
// 1. 解析素材内容
// 2. 添加完整路径到内容中
// 3. 生成素材标题
// 4. 构建并返回响应数据
// 参数：
//
//	medium - 素材实例
//	groupName - 素材分组名称
//
// 返回：素材响应数据
func buildMediumResponseItem(medium model.Medium, groupName string) gin.H {
	// 解析素材内容
	content := service.ParseMediumContent(medium.Content)
	
	// 添加完整路径到内容中
	appendMediumFullPaths(content)
	
	// 生成素材标题
	title := mediumTitle(medium.Type, content)
	
	// 构建并返回响应数据
	return gin.H{
		"id":              medium.ID,              // 素材 ID
		"title":           title,                 // 素材标题
		"type":            service.MediumTypeName(medium.Type), // 素材类型名称
		"typeValue":       medium.Type,           // 素材类型值
		"content":         content,               // 素材内容
		"mediumGroupId":   medium.MediumGroupID,  // 素材分组 ID
		"mediumGroupName": groupName,             // 素材分组名称
		"mediaId":         medium.MediaID,        // 媒体 ID
		"userId":          medium.UserID,         // 用户 ID
		"userName":        medium.UserName,       // 用户名
		"isSync":          medium.IsSync,         // 是否同步
		"createdAt":       medium.CreatedAt.Format("2006-01-02 15:04:05"), // 创建时间
		"updatedAt":       medium.UpdatedAt.Format("2006-01-02 15:04:05"), // 更新时间
	}
}

// appendMediumFullPaths 添加素材完整路径
// 为素材内容中的路径字段添加完整的 URL 路径
// 处理流程：
// 1. 定义需要添加完整路径的字段映射
// 2. 遍历映射，为每个路径字段添加完整 URL
// 参数：
//
//	content - 素材内容
func appendMediumFullPaths(content map[string]interface{}) {
	// 定义需要添加完整路径的字段映射
	pairs := map[string]string{
		"imagePath": "imageFullPath", // 图片路径 -> 图片完整路径
		"videoPath": "videoFullPath", // 视频路径 -> 视频完整路径
		"voicePath": "voiceFullPath", // 语音路径 -> 语音完整路径
		"filePath":  "fileFullPath",  // 文件路径 -> 文件完整路径
	}
	
	// 遍历映射，为每个路径字段添加完整 URL
	for rawKey, fullKey := range pairs {
		if path, ok := content[rawKey].(string); ok && path != "" {
			content[fullKey] = fmt.Sprintf("http://localhost:9501/uploads/%s", path)
		}
	}
}

// mediumTitle 生成素材标题
// 根据素材类型和内容生成素材标题
// 处理流程：
// 1. 根据素材类型，从内容中获取对应的标题字段
// 2. 返回标题字符串
// 参数：
//
//	mediumType - 素材类型
//	content - 素材内容
//
// 返回：素材标题
//
// 类型与标题字段对应关系：
// 1 - 文本：title
// 2 - 图片：imageName
// 3, 6 - 链接/小程序：title
// 4 - 语音：voiceName
// 5 - 视频：videoName
// 7 - 文件：fileName
func mediumTitle(mediumType int, content map[string]interface{}) string {
	switch mediumType {
	case 1:
		return stringValue(content["title"])
	case 2:
		return stringValue(content["imageName"])
	case 3, 6:
		return stringValue(content["title"])
	case 4:
		return stringValue(content["voiceName"])
	case 5:
		return stringValue(content["videoName"])
	case 7:
		return stringValue(content["fileName"])
	default:
		return ""
	}
}

// stringValue 获取字符串值
// 从 interface{} 类型的值中获取字符串值，如果不是字符串类型则返回空字符串
// 处理流程：
// 1. 尝试将值转换为字符串类型
// 2. 如果转换成功，返回字符串值
// 3. 如果转换失败，返回空字符串
// 参数：
//
//	value - 任意类型的值
//
// 返回：字符串值或空字符串
func stringValue(value interface{}) string {
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}

// getUintID 获取无符号整数 ID
// 从 Gin 上下文中获取指定键的无符号整数 ID，支持从路径参数、查询参数和表单参数中获取
// 处理流程：
// 1. 尝试从路径参数中获取 ID
// 2. 如果路径参数中没有，尝试从查询参数中获取
// 3. 如果查询参数中没有，尝试从 "mediumId" 查询参数中获取
// 4. 如果以上都没有，尝试从表单参数中获取
// 5. 如果都没有，返回 0
// 参数：
//
//	c - Gin 上下文
//	key - 键名
//
// 返回：无符号整数 ID 或 0
func getUintID(c *gin.Context, key string) uint64 {
	// 尝试从路径参数中获取 ID
	if val := c.Param(key); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	// 尝试从查询参数中获取 ID
	if val := c.Query(key); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	// 尝试从 "mediumId" 查询参数中获取 ID
	if val := c.Query("mediumId"); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	// 尝试从表单参数中获取 ID
	if val := c.DefaultPostForm(key, ""); val != "" {
		id, _ := strconv.ParseUint(val, 10, 32)
		return id
	}
	// 如果都没有，返回 0
	return 0
}

// asUint 转换为 uint 类型
// 将 interface{} 类型的值转换为 uint 类型，如果转换失败则返回 0
// 处理流程：
// 1. 如果值为 nil，返回 0
// 2. 尝试将值转换为 uint 类型
// 3. 如果转换成功，返回转换后的值
// 4. 如果转换失败，返回 0
// 参数：
//
//	value - 任意类型的值
//
// 返回：uint 类型的值或 0
func asUint(value interface{}) uint {
	// 如果值为 nil，返回 0
	if value == nil {
		return 0
	}
	// 尝试将值转换为 uint 类型
	if v, ok := value.(uint); ok {
		return v
	}
	// 如果转换失败，返回 0
	return 0
}

// MediumGroupHandler 素材分组管理处理器
// 提供素材分组的 CRUD 操作
// 主要职责：
// 1. 获取素材分组列表
// 2. 创建素材分组
// 3. 更新素材分组
// 4. 删除素材分组
//
// 依赖服务：
// - MediumGroupService: 素材分组服务

type MediumGroupHandler struct {
	svc *service.MediumGroupService // 素材分组服务
}

// NewMediumGroupHandler 创建素材分组管理处理器实例
// 参数：db - GORM 数据库连接
// 返回：素材分组管理处理器实例
func NewMediumGroupHandler(db *gorm.DB) *MediumGroupHandler {
	return &MediumGroupHandler{svc: service.NewMediumGroupService(db)}
}

// Index 获取素材分组列表
// 获取指定企业的素材分组列表
// 处理流程：
// 1. 获取企业 ID
// 2. 调用服务层获取分组列表
// 3. 返回分组列表
//
// 返回：包含素材分组列表的响应
func (h *MediumGroupHandler) Index(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}
	
	// 调用服务层获取分组列表
	groups, err := h.svc.List(corpID.(uint))
	if err != nil {
		response.Fail(c, response.ErrDB, "获取素材分组列表失败")
		return
	}
	
	// 返回分组列表
	response.Success(c, groups)
}

// Store 创建素材分组
// 创建新的素材分组
// 处理流程：
// 1. 绑定请求参数
// 2. 调用服务层创建分组
// 3. 返回创建的分组
// 请求体（JSON）：
//
//	包含素材分组的相关信息，如名称、企业 ID 等
//
// 返回：创建的素材分组详情
func (h *MediumGroupHandler) Store(c *gin.Context) {
	// 绑定请求参数
	var group model.MediumGroup
	if err := c.ShouldBindJSON(&group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层创建分组
	if err := h.svc.Create(&group); err != nil {
		response.Fail(c, response.ErrDB, "创建素材分组失败")
		return
	}
	
	// 返回创建的分组
	response.Success(c, group)
}

// Update 更新素材分组
// 更新现有素材分组的信息
// 处理流程：
// 1. 获取分组 ID
// 2. 创建分组实例并设置 ID
// 3. 绑定请求参数
// 4. 调用服务层更新分组
// 5. 返回更新后的分组
// 参数：
//
//	id - 分组 ID（路径参数）
//
// 请求体（JSON）：
//
//	包含素材分组的更新信息，如名称等
//
// 返回：更新后的素材分组详情
func (h *MediumGroupHandler) Update(c *gin.Context) {
	// 获取分组 ID
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	// 创建分组实例并设置 ID
	group := &model.MediumGroup{}
	group.ID = uint(id)
	
	// 绑定请求参数
	if err := c.ShouldBindJSON(group); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	
	// 调用服务层更新分组
	if err := h.svc.Update(group); err != nil {
		response.Fail(c, response.ErrDB, "更新素材分组失败")
		return
	}
	
	// 返回更新后的分组
	response.Success(c, group)
}

// Destroy 删除素材分组
// 删除指定的素材分组
// 处理流程：
// 1. 获取分组 ID
// 2. 调用服务层删除分组
// 3. 返回删除结果
// 参数：
//
//	id - 分组 ID（路径参数）
//
// 返回：删除成功的消息
func (h *MediumGroupHandler) Destroy(c *gin.Context) {
	// 获取分组 ID
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	// 调用服务层删除分组
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除素材分组失败")
		return
	}
	
	// 返回删除结果
	response.SuccessMsg(c, "删除成功")
}
