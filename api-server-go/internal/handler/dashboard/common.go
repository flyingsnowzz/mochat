// Package dashboard 提供 Dashboard 相关的 HTTP 处理器
// 该包包含通用功能、欢迎语管理、聊天工具等处理逻辑
// 主要职责：
// 1. 处理文件上传等通用功能
// 2. 管理企业微信应用
// 3. 处理欢迎语的增删改查
// 4. 提供聊天工具栏管理
package dashboard

import (
	"encoding/json"
	"strconv"
	"strings"

	"mochat-api-server/internal/model"
	"mochat-api-server/internal/pkg/response"
	"mochat-api-server/internal/pkg/storage"
	"mochat-api-server/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommonHandler 通用功能处理器
// 提供文件上传等通用功能
type CommonHandler struct{}

// NewCommonHandler 创建通用功能处理器实例
func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

// Upload 文件上传处理函数
// 接收表单文件并上传到存储系统，返回文件访问 URL
// 处理流程：
// 1. 从请求中获取上传的文件
// 2. 生成唯一的文件名
// 3. 打开文件并读取内容
// 4. 调用存储服务上传文件
// 5. 生成文件访问 URL 并返回
// 参数：c - Gin 上下文
// 返回：包含文件 URL、路径、名称、大小的成功响应
func (h *CommonHandler) Upload(c *gin.Context) {
	// 从表单中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, response.ErrParams, "请上传文件")
		return
	}

	// 生成唯一的文件名，避免文件冲突
	filename := storage.GenerateFilename(file.Filename)

	// 打开文件准备上传
	src, err := file.Open()
	if err != nil {
		response.Fail(c, response.ErrFileUpload, "读取文件失败")
		return
	}
	defer src.Close() // 确保文件在处理完成后关闭

	// 调用存储服务上传文件
	path, err := storage.DefaultStorage.Upload(src, filename)
	if err != nil {
		response.Fail(c, response.ErrFileUpload, "上传文件失败")
		return
	}

	// 生成文件访问 URL
	url := storage.DefaultStorage.GetURL(path)

	// 返回成功响应，包含文件相关信息
	response.Success(c, gin.H{
		"url":  url,           // 文件访问 URL
		"path": path,          // 文件存储路径
		"name": file.Filename, // 原始文件名
		"size": file.Size,     // 文件大小
	})
}

// UploadFile 文件上传处理函数（Upload 方法的别名
func (h *CommonHandler) UploadFile(c *gin.Context) {
	h.Upload(c)
}

// IndexHandler 首页数据处理器
// 处理首页相关的数据请求
type IndexHandler struct {
	dataSvc *service.CorpDayDataService
}

// NewIndexHandler 创建首页数据处理器实例
// 参数：db - GORM 数据库连接
func NewIndexHandler(db *gorm.DB) *IndexHandler {
	return &IndexHandler{dataSvc: service.NewCorpDayDataService(db)}
}

// Index 首页数据接口
// 返回首页相关数据（当前未实现完整功能）
func (h *IndexHandler) Index(c *gin.Context) {
	response.Success(c, gin.H{})
}

// LineChat 折线图数据接口
// 返回折线图相关数据（当前未实现完整功能）
func (h *IndexHandler) LineChat(c *gin.Context) {
	response.Success(c, gin.H{})
}

// ChatToolHandler 聊天工具栏处理器
// 管理侧边栏聊天工具相关功能
type ChatToolHandler struct {
	svc *service.ChatToolService
}

// NewChatToolHandler 创建聊天工具栏处理器实例
// 参数：db - GORM 数据库连接
func NewChatToolHandler(db *gorm.DB) *ChatToolHandler {
	return &ChatToolHandler{svc: service.NewChatToolService(db)}
}

// Index 获取聊天工具栏列表
// 返回所有可用的聊天工具列表
func (h *ChatToolHandler) Index(c *gin.Context) {
	tools, err := h.svc.List()
	if err != nil {
		response.Fail(c, response.ErrDB, "获取侧边工具栏列表失败")
		return
	}
	response.Success(c, tools)
}

// AgentHandler 应用管理处理器
// 处理企业微信应用相关操作
type AgentHandler struct {
	svc *service.WorkAgentService
}

// NewAgentHandler 创建应用管理处理器实例
// 参数：db - GORM 数据库连接
func NewAgentHandler(db *gorm.DB) *AgentHandler {
	return &AgentHandler{svc: service.NewWorkAgentService(db)}
}

// Store 创建新应用
// 接收应用信息并保存到数据库
func (h *AgentHandler) Store(c *gin.Context) {
	var agent model.WorkAgent
	if err := c.ShouldBindJSON(&agent); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Create(&agent); err != nil {
		response.Fail(c, response.ErrDB, "创建应用失败")
		return
	}
	response.Success(c, agent)
}

// TxtVerifyShow 文本验证展示接口
func (h *AgentHandler) TxtVerifyShow(c *gin.Context) {
	response.Success(c, gin.H{})
}

// TxtVerifyUpload 文本验证文件上传接口
func (h *AgentHandler) TxtVerifyUpload(c *gin.Context) {
	response.SuccessMsg(c, "上传成功")
}

// GreetingHandler 欢迎语管理处理器
// 处理客户欢迎语的增删改查操作
type GreetingHandler struct {
	db  *gorm.DB
	svc *service.GreetingService
}

// NewGreetingHandler 创建欢迎语管理处理器实例
// 参数：db - GORM 数据库连接
func NewGreetingHandler(db *gorm.DB) *GreetingHandler {
	return &GreetingHandler{db: db, svc: service.NewGreetingService(db)}
}

// Index 获取欢迎语列表
// 支持分页查询企业的欢迎语列表，返回包含员工信息和素材内容
// 处理流程：
// 1. 获取企业 ID 和分页参数
// 2. 调用服务层获取欢迎语列表
// 3. 构建返回数据，包括：
//   - 欢迎语列表项
//   - 是否存在通用欢迎语
//   - 已配置欢迎语的员工列表
//   - 分页信息
//
// 参数：
//
//	page - 页码，默认为 1
//	perPage/pageSize - 每页数量，默认为 10
//
// 返回：包含欢迎语列表、是否有通用欢迎语、已配置的员工列表、分页信息的响应
func (h *GreetingHandler) Index(c *gin.Context) {
	// 获取企业 ID
	corpID, _ := c.Get("corpId")
	if corpID == nil {
		corpID = uint(0)
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", c.DefaultQuery("pageSize", "10")))

	// 调用服务层获取欢迎语列表
	greetings, total, err := h.svc.List(corpID.(uint), page, pageSize)
	if err != nil {
		response.Fail(c, response.ErrDB, "获取欢迎语列表失败")
		return
	}

	// 构建返回数据
	items := make([]gin.H, 0, len(greetings))
	hadGeneral := 0                         // 是否存在通用欢迎语（1-存在，0-不存在）
	hadEmployees := make([]uint, 0)         // 已配置欢迎语的员工 ID 列表
	seenEmployee := make(map[uint]struct{}) // 用于去重员工 ID

	// 处理每个欢迎语
	for _, greeting := range greetings {
		// 构建欢迎语列表项
		item := h.buildGreetingListItem(greeting)
		items = append(items, item)

		// 检查是否为通用欢迎语
		if greeting.RangeType == 1 {
			hadGeneral = 1
		}

		// 收集已配置欢迎语的员工 ID
		for _, employee := range greetingEmployeeIDs(greeting.Employees) {
			if _, ok := seenEmployee[employee]; ok {
				continue // 跳过已处理的员工
			}
			seenEmployee[employee] = struct{}{}
			hadEmployees = append(hadEmployees, employee)
		}
	}

	// 返回成功响应
	response.Success(c, gin.H{
		"list":         items,        // 欢迎语列表
		"hadGeneral":   hadGeneral,   // 是否存在通用欢迎语
		"hadEmployees": hadEmployees, // 已配置欢迎语的员工列表
		"page": gin.H{
			"perPage":     pageSize,                       // 每页数量
			"total":       total,                          // 总记录数
			"totalPage":   calcTotalPage(total, pageSize), // 总页数
			"currentPage": page,                           // 当前页码
		},
		"currentPage": page,     // 当前页码（兼容旧接口）
		"pageSize":    pageSize, // 每页数量（兼容旧接口）
		"total":       total,    // 总记录数（兼容旧接口）
	})
}

// Show 获取欢迎语详情
// 根据 ID 获取指定欢迎语的详细信息
// 参数：id/greetingId - 欢迎语 ID
// 返回：欢迎语详情信息
func (h *GreetingHandler) Show(c *gin.Context) {
	id := getGreetingUintID(c, "id", "greetingId")
	greeting, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "欢迎语不存在")
		return
	}
	response.Success(c, h.buildGreetingDetailItem(*greeting))
}

// Store 创建新欢迎语
// 创建新的客户欢迎语配置
// 参数（JSON Body）：
//
//	rangeType - 范围类型（1-全部，2-指定）
//	employees - 员工 ID 列表
//	type - 欢迎语类型
//	words - 欢迎语文本
//	mediumId - 素材 ID
//
// 返回：创建的欢迎语详情
func (h *GreetingHandler) Store(c *gin.Context) {
	var req struct {
		RangeType int    `json:"rangeType" binding:"required"`
		Employees string `json:"employees"`
		Type      string `json:"type" binding:"required"`
		Words     string `json:"words"`
		MediumID  uint   `json:"mediumId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	corpID, _ := c.Get("corpId")
	if corpID == nil || corpID.(uint) == 0 {
		response.Fail(c, response.ErrParams, "请先选择企业")
		return
	}
	greeting := model.Greeting{
		CorpID:    corpID.(uint),
		Type:      strings.TrimSpace(req.Type),
		Words:     req.Words,
		MediumID:  req.MediumID,
		RangeType: req.RangeType,
		Employees: greetingEmployeeJSON(req.Employees),
	}
	if err := h.svc.Create(&greeting); err != nil {
		response.Fail(c, response.ErrDB, "创建欢迎语失败")
		return
	}
	response.Success(c, h.buildGreetingDetailItem(greeting))
}

// Update 更新欢迎语
// 更新指定欢迎语的配置信息
// 参数（JSON Body）：
//
//	greetingId - 欢迎语 ID
//	rangeType - 范围类型
//	employees - 员工 ID 列表
//	type - 欢迎语类型
//	words - 欢迎语文本
//	mediumId - 素材 ID
//
// 返回：更新后的欢迎语详情
func (h *GreetingHandler) Update(c *gin.Context) {
	id := getGreetingUintID(c, "id", "greetingId")
	var greeting *model.Greeting
	var err error

	var req struct {
		GreetingID uint   `json:"greetingId"`
		RangeType  int    `json:"rangeType" binding:"required"`
		Employees  string `json:"employees"`
		Type       string `json:"type" binding:"required"`
		Words      string `json:"words"`
		MediumID   uint   `json:"mediumId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if id == 0 {
		id = uint64(req.GreetingID)
		if id == 0 {
			response.Fail(c, response.ErrParams, "参数错误")
			return
		}
	}
	greeting, err = h.svc.GetByID(uint(id))
	if err != nil {
		response.Fail(c, response.ErrNotFound, "欢迎语不存在")
		return
	}
	greeting.RangeType = req.RangeType
	greeting.Employees = greetingEmployeeJSON(req.Employees)
	greeting.Type = strings.TrimSpace(req.Type)
	greeting.Words = req.Words
	greeting.MediumID = req.MediumID
	if err := h.svc.Update(greeting); err != nil {
		response.Fail(c, response.ErrDB, "更新欢迎语失败")
		return
	}
	response.Success(c, h.buildGreetingDetailItem(*greeting))
}

// Destroy 删除欢迎语
// 根据 ID 删除指定的欢迎语
// 参数：id/greetingId - 欢迎语 ID
func (h *GreetingHandler) Destroy(c *gin.Context) {
	id := getGreetingUintID(c, "id", "greetingId")
	if id == 0 {
		var req struct {
			GreetingID uint `json:"greetingId"`
		}
		if err := c.ShouldBindJSON(&req); err == nil {
			id = uint64(req.GreetingID)
		}
	}
	if id == 0 {
		response.Fail(c, response.ErrParams, "参数错误")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		response.Fail(c, response.ErrDB, "删除欢迎语失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// buildGreetingListItem 构建欢迎语列表项
// 将模型数据转换为前端需要的格式
// 参数：greeting - 欢迎语模型
// 返回：包含列表项数据
func (h *GreetingHandler) buildGreetingListItem(greeting model.Greeting) gin.H {
	employeeIDs := greetingEmployeeIDs(greeting.Employees)
	employeeMap := h.loadGreetingEmployees(employeeIDs)
	employeeNames := make([]string, 0, len(employeeIDs))
	for _, employeeID := range employeeIDs {
		employee, ok := employeeMap[employeeID]
		if !ok {
			continue
		}
		employeeNames = append(employeeNames, employee.Name)
	}

	mediumContent := gin.H{}
	typeText := greetingTypeText(greeting.Type)
	if greeting.MediumID > 0 {
		if medium, err := service.NewMediumService(h.db).GetByID(greeting.MediumID); err == nil {
			content := service.ParseMediumContent(medium.Content)
			appendMediumFullPaths(content)
			mediumContent = gin.H(content)
			if greeting.Words == "" {
				typeText = service.MediumTypeName(medium.Type)
			}
		}
	}

	return gin.H{
		"greetingId":    greeting.ID,
		"rangeType":     greeting.RangeType,
		"rangeTypeText": greetingRangeTypeText(greeting.RangeType),
		"employees":     employeeNames,
		"words":         greeting.Words,
		"mediumId":      greeting.MediumID,
		"mediumContent": mediumContent,
		"type":          greeting.Type,
		"typeText":      typeText,
		"createdAt":     greeting.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// buildGreetingDetailItem 构建欢迎语详情项
// 将模型数据转换为前端需要的详情格式
// 参数：greeting - 欢迎语模型
// 返回：包含详情数据
func (h *GreetingHandler) buildGreetingDetailItem(greeting model.Greeting) gin.H {
	employeeIDs := greetingEmployeeIDs(greeting.Employees)
	employeeMap := h.loadGreetingEmployees(employeeIDs)
	employees := make([]gin.H, 0, len(employeeIDs))
	for _, employeeID := range employeeIDs {
		employee, ok := employeeMap[employeeID]
		if !ok {
			continue
		}
		employees = append(employees, gin.H{
			"employeeId":   employee.ID,
			"employeeName": employee.Name,
			"name":         employee.Name,
			"id":           employee.ID,
		})
	}

	mediumContent := gin.H{}
	typeText := greetingTypeText(greeting.Type)
	if greeting.MediumID > 0 {
		if medium, err := service.NewMediumService(h.db).GetByID(greeting.MediumID); err == nil {
			content := service.ParseMediumContent(medium.Content)
			appendMediumFullPaths(content)
			mediumContent = gin.H(content)
			if greeting.Words == "" {
				typeText = service.MediumTypeName(medium.Type)
			}
		}
	}

	return gin.H{
		"greetingId":    greeting.ID,
		"rangeType":     greeting.RangeType,
		"rangeTypeText": greetingRangeTypeText(greeting.RangeType),
		"employees":     employees,
		"words":         greeting.Words,
		"mediumId":      greeting.MediumID,
		"mediumContent": mediumContent,
		"type":          greeting.Type,
		"typeText":      typeText,
		"createdAt":     greeting.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// loadGreetingEmployees 加载欢迎语关联的员工信息
// 根据员工 ID 列表批量加载员工数据
// 参数：ids - 员工 ID 列表
// 返回：员工 ID 到员工模型的映射
func (h *GreetingHandler) loadGreetingEmployees(ids []uint) map[uint]model.WorkEmployee {
	result := make(map[uint]model.WorkEmployee)
	if len(ids) == 0 {
		return result
	}
	var employees []model.WorkEmployee
	if err := h.db.Where("id IN ?", ids).Find(&employees).Error; err != nil {
		return result
	}
	for _, employee := range employees {
		result[employee.ID] = employee
	}
	return result
}

// greetingEmployeeJSON 将员工数据转换为 JSON 字符串
// 参数：raw - 原始员工数据（JSON 字符串或逗号分隔字符串）
// 返回：标准 JSON 数组字符串
func greetingEmployeeJSON(raw string) string {
	ids := greetingEmployeeIDs(raw)
	if len(ids) == 0 {
		return "[]"
	}
	data, err := json.Marshal(ids)
	if err != nil {
		return "[]"
	}
	return string(data)
}

// greetingEmployeeIDs 解析员工 ID 列表
// 支持 JSON 数组和逗号分隔两种格式
// 参数：raw - 原始员工数据
// 返回：员工 ID 列表
func greetingEmployeeIDs(raw string) []uint {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []uint{}
	}
	if strings.HasPrefix(raw, "[") {
		var ids []uint
		if err := json.Unmarshal([]byte(raw), &ids); err == nil {
			return ids
		}
	}

	parts := strings.Split(raw, ",")
	result := make([]uint, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			continue
		}
		result = append(result, uint(id))
	}
	return result
}

// greetingTypeText 转换欢迎语类型为文本
// 参数：greetingType - 欢迎语类型字符串
// 返回：类型名称文本
func greetingTypeText(greetingType string) string {
	parts := strings.Split(greetingType, ",")
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch part {
		case "1":
			names = append(names, "文本")
		case "2":
			names = append(names, "图片")
		case "3":
			names = append(names, "图文")
		case "6":
			names = append(names, "小程序")
		}
	}
	return strings.Join(names, "+")
}

// greetingRangeTypeText 转换范围类型为文本
// 参数：rangeType - 范围类型（1-全部，2-指定）
// 返回：范围类型文本
func greetingRangeTypeText(rangeType int) string {
	if rangeType == 2 {
		return "指定成员"
	}
	return "全部成员"
}

// calcTotalPage 计算总页数
// 参数：
//
//	total - 总记录数
//	pageSize - 每页数量
//
// 返回：总页数
func calcTotalPage(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	if total == 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}

// getGreetingUintID 获取欢迎语 ID
// 从路径参数或查询参数中获取 ID，支持多个参数名
// 参数：
//
//	c - Gin 上下文
//	keys - 可能的参数名列表
//
// 返回：ID 值
func getGreetingUintID(c *gin.Context, keys ...string) uint64 {
	for _, key := range keys {
		if val := c.Param(key); val != "" {
			id, _ := strconv.ParseUint(val, 10, 32)
			return id
		}
		if val := c.Query(key); val != "" {
			id, _ := strconv.ParseUint(val, 10, 32)
			return id
		}
	}
	return 0
}
